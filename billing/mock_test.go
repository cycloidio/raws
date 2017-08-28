package billing

import (
	"io"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cycloidio/raws"
)

type mockChecker struct {
	Checker

	// Check returned elements
	cb bool
	ce error

	// AlreadyPresent returned elements
	apb bool
	aps string
}

func (m mockChecker) Check(bucket string, filename string) (bool, error) {
	return m.cb, m.ce
}

func (m mockChecker) AlreadyPresent() (bool, string) {
	return m.apb, m.aps
}

type mockDownloader struct {
	Downloader

	// Download returned elements
	ds string
	de error

	// Unzip returned elements
	us string
	ue error
}

func (m mockDownloader) Download(bucket string, filename string, dest string) (string, error) {
	return m.ds, m.de
}

func (m mockDownloader) Unzip(src string, dest string) (string, error) {
	return m.us, m.ue
}

type mockLoader struct {
	Loader

	// ProcessFile returned elements
	pfs []string
	pfe error

	// GetStats returned elements
	gss *stats
}

func (m mockLoader) ProcessFile(reportName string, billingFile string) ([]string, error) {
	return m.pfs, m.pfe
}

func (m mockLoader) TerminateProcessFile() {
	return
}

func (m mockLoader) GetStats() *stats {
	return m.gss
}

type mockInjector struct {
	Injector

	// CreateRecord returned element
	crde error

	// CreateReport returned element
	crpe error

	// CreateRecords returned elements
	crss []string
	crsi int
	crse error

	// MaxRecords returned element
	mri int
}

func (m mockInjector) CreateRecord(record *billingRecord) error {
	return m.crde
}

func (m mockInjector) CreateRecords(records []*billingRecord) ([]string, int, error) {
	return m.crss, m.crsi, m.crse
}

func (m mockInjector) CreateReport(filename string, hash string) error {
	return m.crpe
}

func (m mockInjector) MaxRecords() int {
	return m.mri
}

// mockAWSReader used for testing purposes
type mockAWSReader struct {
	raws.AWSReader

	// GetAccountID returned element
	accountId string

	// ListObjects returned elements
	loo []*s3.ListObjectsOutput
	loe raws.Errs

	// DownloadObject returned elements
	doi int64
	doe error
}

func (m mockAWSReader) GetAccountID() string {
	return m.accountId
}

func (m mockAWSReader) ListObjects(input *s3.ListObjectsInput) ([]*s3.ListObjectsOutput, raws.Errs) {
	return m.loo, m.loe
}

func (m mockAWSReader) DownloadObject(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error) {
	return m.doi, m.doe
}

func checkErrors(t *testing.T, name string, index int, err error, expected error) {
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("%s [%d] - error: received=%+v | expected=%+v",
			name, index, err, expected)
	}
}

type mockDynamodb struct {
	dynamodbiface.DynamoDBAPI

	// GetItem returned elements
	geo *dynamodb.GetItemOutput
	gee error

	// PutItem returned elements
	pio *dynamodb.PutItemOutput
	pie error

	// BatchWriteItem returned elements
	bwio *dynamodb.BatchWriteItemOutput
	bwie error
}

func (m mockDynamodb) GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return m.geo, m.gee
}

func (m mockDynamodb) PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return m.pio, m.pie
}

func (m mockDynamodb) BatchWriteItem(*dynamodb.BatchWriteItemInput) (*dynamodb.BatchWriteItemOutput, error) {
	return m.bwio, m.bwie
}
