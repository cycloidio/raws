package billing

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/twinj/uuid"
)

type Loader interface {
	ProcessFile(reportName string, billingFile string) ([]string, error)
	TerminateProcessFile()
	GetStats() *stats
}

func NewLoader(injector Injector) Loader {
	return &billingLoader{
		json:     []byte{},
		injector: injector,
		result:   newStats(),
		reportFd: nil,
	}
}

func (l *billingLoader) ProcessFile(reportName string, billingFile string) ([]string, error) {
	var openErr error
	var end = false

	if l.reportFd == nil {
		l.reportName = reportName
		l.report, openErr = l.openBillingReport(billingFile)
		if openErr != nil {
			return nil, openErr
		}
	}
	for end == false {
		l.result.read++
		record := &billingRecord{}
		values, readErr := l.report.csvReader.Read()
		if readErr == io.EOF {
			end = true
		} else if readErr != nil {
			l.result.warnings++
			return nil, NewCSVError(readErr)
		}
		if parseErr := l.parseRecord(values, record, l.report); parseErr != nil {
			l.result.warnings++
			return nil, NewConvertError(parseErr)
		}
		if recordIds, saveErr := l.saveRecords(record, end); saveErr != nil {
			return recordIds, saveErr
		}
	}
	return nil, nil
}

func (l *billingLoader) TerminateProcessFile() {
	l.reportFd.Close()
}

func (l *billingLoader) GetStats() *stats {
	return l.result
}

func newStats() *stats {
	return &stats{
		read:     1,
		loaded:   0,
		warnings: 0,
		failed:   0,
	}
}

type stats struct {
	read     int
	loaded   int
	warnings int
	failed   int
}

type billingLoader struct {
	reportName string
	report     *billingReport
	records    []*billingRecord
	reportFd   *os.File
	json       []byte
	injector   Injector
	result     *stats
}

type billingReport struct {
	CsvFileName   string
	InvoicePeriod time.Time
	Fields        []string
	csvReader     *csv.Reader
}

type billingRecord struct {
	Id               string
	ReportName       string
	InvoiceID        string
	PayerAccountId   uint64
	LinkedAccountId  uint64
	RecordType       string
	RecordId         string
	ProductName      string
	RateId           int64
	SubscriptionId   int64
	PricingPlanId    int64
	UsageType        string
	Operation        string
	AvailabilityZone string
	ReservedInstance string
	ItemDescription  string
	UsageStartDate   string
	UsageEndDate     string
	UsageQuantity    float64
	BlendedRate      float64
	BlendedCost      float64
	UnBlendedRate    float64
	UnBlendedCost    float64
	ResourceId       string
	Tags             map[string]string
}

func (l *billingLoader) parseRecord(in []string, out *billingRecord, report *billingReport) error {
	var tagMatcher = regexp.MustCompile("(user|aws):.*")
	var fieldTypes = map[string]func(s string, report *billingReport) (interface{}, error){
		"PayerAccountId":  l.parseUint,
		"LinkedAccountId": l.parseUint,
		"RecordInt":       l.parseUint,
		"RateId":          l.parseInt,
		"SubscriptionId":  l.parseInt,
		"PricingPlanId":   l.parseInt,
		"UsageStartDate":  l.parseDate,
		"UsageEndDate":    l.parseDate,
		"UsageQuantity":   l.parseFloat,
		"BlendedRate":     l.parseFloat,
		"BlendedCost":     l.parseFloat,
		"UnBlendedRate":   l.parseFloat,
		"UnBlendedCost":   l.parseFloat,
	}

	if len(in) == 0 {
		return nil
	}
	out.Tags = make(map[string]string)
	for i, field := range report.Fields {
		if f, ok := fieldTypes[field]; ok {
			value, err := f(in[i], report)
			if err != nil {
				return err
			}
			err = l.assignField(out, field, value, false)
			if err != nil {
				return err
			}
		} else if tagMatcher.MatchString(field) {
			key, err := l.parseTag(field, report)
			if err != nil {
				return err
			}
			value := map[string]string{key.(string): in[i]}
			if err := l.assignField(out, field, value, true); err != nil {
				return err
			}
		} else {
			if err := l.assignField(out, field, in[i], false); err != nil {
				return err
			}
		}
	}
	if out.RecordId == "0" || out.RecordId == "" {
		return fmt.Errorf("no recordId found for this entry: %v", in)
	}
	return nil
}

func (l *billingLoader) assignField(record *billingRecord, field string, value interface{}, tag bool) error {
	if tag == false {
		recordField := reflect.ValueOf(record).Elem().FieldByName(field)
		if !recordField.IsValid() {
			return fmt.Errorf("no such field '%s'", field)
		}
		if !recordField.CanSet() {
			return fmt.Errorf("field '%s' cannot be set", field)
		}
		valueTyped := reflect.ValueOf(value)
		if recordField.Type() != valueTyped.Type() {
			return fmt.Errorf("field '%s' didn't have type %s", field, valueTyped.Type().Name())
		}
		recordField.Set(valueTyped)
	} else {
		for k, v := range value.(map[string]string) {
			record.Tags[k] = v
		}
	}
	return nil
}

func (l *billingLoader) openBillingReport(billingFile string) (*billingReport, error) {
	var dateMonthMatcher = regexp.MustCompile(`\d{4}-\d{2}`)

	invoicePeriod, err := time.Parse("2006-02", dateMonthMatcher.FindString(billingFile))
	if err != nil {
		return nil, err
	}
	file, err := os.Open(billingFile)
	if err != nil {
		return nil, err
	}
	l.reportFd = file
	reader := csv.NewReader(file)
	report := &billingReport{
		CsvFileName:   billingFile,
		InvoicePeriod: invoicePeriod,
		csvReader:     reader}

	fields, err := reader.Read()
	if err != nil {
		return nil, err
	}
	report.Fields = fields
	return report, nil
}

func (l *billingLoader) parseTag(s string, report *billingReport) (interface{}, error) {
	tagParts := strings.Split(s, ":")
	return strings.Join(tagParts, "_"), nil
}

func (l *billingLoader) parseInt(s string, report *billingReport) (interface{}, error) {
	return strconv.ParseInt(s, 0, 0)
}

func (l *billingLoader) parseUint(s string, report *billingReport) (interface{}, error) {
	return strconv.ParseUint(s, 0, 0)
}

func (l *billingLoader) parseFloat(s string, report *billingReport) (interface{}, error) {
	return strconv.ParseFloat(s, 0)
}

func (l *billingLoader) parseDate(s string, report *billingReport) (interface{}, error) {
	var returnTime time.Time
	var err error
	switch s {
	case "":
		returnTime = report.InvoicePeriod
	default:
		returnTime, err = time.Parse("2006-01-02 15:04:05", s)
		if err != nil {
			return nil, err
		}
	}
	return returnTime.Format(time.RFC3339), nil
}

func (l *billingLoader) saveRecords(in *billingRecord, flush bool) ([]string, error) {
	if len(l.records) < l.injector.MaxRecords() {
		in.ReportName = l.reportName
		in.Id = uuid.NewV4().String()
		l.records = append(l.records, in)
	}
	if len(l.records) == l.injector.MaxRecords() || flush == true {
		recordIds, processed, err := l.injector.CreateRecords(l.records)
		if err != nil {
			log.Debugf("Loader - failed loading: %v - %v", recordIds, err)
		} else {
			log.Debugf("Loader - Successfully loaded: %d items", processed)
		}
		l.records = nil
		l.result.failed += len(recordIds)
		l.result.loaded += processed
		return recordIds, err
	}
	return nil, nil
}
