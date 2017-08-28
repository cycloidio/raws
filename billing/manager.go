package billing

import (
	"os"

	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/cycloidio/raws"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors:    false,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

type Manager interface {
	ImportFromS3(date string, bucket string) error
	ImportFromFile(reportName string, filepath string) error
	SetLogger(textFormat *log.TextFormatter, file *io.Writer, level *log.Level)
}

type BillingManager struct {
	date          string
	s3Connector   raws.AWSReader
	dynamoSvc     *dynamodb.DynamoDB
	dynamoAccount *AwsConfig
	s3Account     *AwsConfig
	checker       Checker
	downloader    Downloader
	loader        Loader
	injector      Injector
}

type AwsConfig struct {
	AccessKey string
	SecretKey string
	Region    string
}

func NewManager(dynamoAccount *AwsConfig, s3Account *AwsConfig) (Manager, error) {
	connector, err := raws.NewAWSReader(
		s3Account.AccessKey,
		s3Account.SecretKey,
		[]string{s3Account.Region},
		nil,
	)
	if err != nil {
		return nil, err
	}
	svc, err := initDynamoService(dynamoAccount)
	if err != nil {
		return nil, err
	}
	injector := NewInjector(svc)
	return &BillingManager{
		s3Connector:   connector,
		s3Account:     s3Account,
		dynamoAccount: dynamoAccount,
		dynamoSvc:     svc,
		checker:       NewChecker(connector, svc),
		downloader:    NewDownloader(connector),
		injector:      injector,
		loader:        NewLoader(injector),
	}, nil
}

func (m *BillingManager) SetLogger(textFormat *log.TextFormatter, file *io.Writer, level *log.Level) {
	if textFormat != nil {
		log.SetFormatter(textFormat)
	}
	if file != nil {
		log.SetOutput(*file)
	}
	if level != nil {
		log.SetLevel(*level)
	}
}

func (m *BillingManager) ImportFromS3(date string, bucket string) error {
	const (
		downloadPath = "/tmp/billing-reports-download/"
		unzipPath    = "/tmp/billing-reports-unzip/"
	)
	var processErr error
	var processRecordIds []string
	m.date = date

	log.Infof("Manager - starting to import %s from S3...", m.getS3Filename())
	defer log.Info("Manager - import from S3 done.")
	needImport, err := m.checker.Check(bucket, m.getS3Filename())
	if err != nil {
		log.Errorf("Checker - error during check: %v", err)
		return err
	}
	if needImport == false {
		log.Info("Checker - file doesn't need import.")
		return nil
	}
	log.Info("Checker - file needs import.")
	downloadedFile, err := m.downloader.Download(bucket, m.getS3Filename(), downloadPath)
	if err != nil {
		log.Errorf("Downloader - error during download: %v", err)
		return err
	}
	log.Info("Downloader - file succesfuly downloaded.")
	filePath, err := m.downloader.Unzip(downloadedFile, unzipPath)
	if err != nil {
		log.Errorf("Downloader - file couldn't be unzipped: %v", err)
		return err
	}
	log.Info("Downloader - file succesfuly unzipped.")
	log.Info("Loader - starting to import file (might take a while)...")
	for processErr == nil {
		processRecordIds, processErr = m.loader.ProcessFile(m.getS3Filename(), filePath)
		if processErr == nil {
			break
		}
		if IsDynamoDBError(processErr) {
			log.Errorf("Loader - couldn't inject following records:")
			for _, recordId := range processRecordIds {
				log.Errorf("Loader - FAILED: %s - %s", recordId, m.getS3Filename())
			}
		} else if IsConvertError(processErr) {
			log.Warningf("Loader - conversion issue: %v", processErr)
		} else if IsCSVError(processErr) {
			log.Warningf("Loader - reading CSV issue: %v", processErr)
		} else {
			log.Errorf("Loader - cannot import file: %v", processErr)
			return processErr
		}
		processErr = nil
		processRecordIds = nil
	}
	m.loader.TerminateProcessFile()
	log.Info("Loader - ...done!")
	present, hash := m.checker.AlreadyPresent()
	if present {
		return nil
	}
	log.Info("Injector - creating report entry...")
	err = m.injector.CreateReport(m.getS3Filename(), hash)
	if err != nil {
		log.Errorf("Injector - error during entry creation: %v", err)
		return err
	} else {
		log.Info("Injector - ...done!")
	}
	stats := m.loader.GetStats()
	log.Infof("Manager - loaded %d/%d: warnings: %d, failed: %d",
		stats.loaded, stats.read,
		stats.warnings, stats.failed)
	return nil
}

func (m *BillingManager) ImportFromFile(reportName string, filePath string) error {
	var processErr error
	var processRecordIds []string

	log.Infof("Manager - starting to import %s from local file...", reportName)
	defer log.Infof("Manager - import from local file done.")
	log.Info("Loader - starting to import file (might take a while)...")
	for {
		processRecordIds, processErr = m.loader.ProcessFile(m.getS3Filename(), filePath)
		if processErr == nil {
			break
		}
		if IsDynamoDBError(processErr) {
			log.Errorf("Loader - couldn't inject following records:")
			for _, recordId := range processRecordIds {
				log.Errorf("Loader - FAILED: %s - %s", recordId, m.getS3Filename())
			}
		} else if IsConvertError(processErr) {
			log.Warningf("Loader - conversion issue: %v", processErr)
		} else if IsCSVError(processErr) {
			log.Warningf("Loader - reading CSV issue: %v", processErr)
		} else {
			log.Errorf("Loader - cannot import file: %v", processErr)
			return processErr
		}
		processErr = nil
		processRecordIds = nil
	}
	m.loader.TerminateProcessFile()
	log.Info("Loader - ...done!")
	stats := m.loader.GetStats()
	log.Infof("Manager - loaded %d/%d: warnings: %d, failed: %d",
		stats.loaded, stats.read,
		stats.warnings, stats.failed)
	return nil
}

func initDynamoService(config *AwsConfig) (*dynamodb.DynamoDB, error) {
	var token string

	creds := credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, token)
	_, err := creds.Get()
	if err != nil {
		return nil, err
	}
	session := session.Must(
		session.NewSession(&aws.Config{
			Region:      aws.String(config.Region),
			DisableSSL:  aws.Bool(false),
			MaxRetries:  aws.Int(3),
			Credentials: creds,
		}),
	)
	return dynamodb.New(session), nil
}

func (m *BillingManager) getS3Filename() string {
	const (
		filenamePattern = "-aws-billing-detailed-line-items-with-resources-and-tags-"
		fileExtension   = ".csv.zip"
	)
	return m.s3Connector.GetAccountID() + filenamePattern + m.date + fileExtension
}
