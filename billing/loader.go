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

	"encoding/json"

	"github.com/twinj/uuid"
)

type Loader interface {
	ProcessFile(reportName string, billingFile string) error
}

func NewLoader(injector Injector) Loader {
	return &billingLoader{
		json:     []byte{},
		injector: injector,
	}
}

func (l *billingLoader) ProcessFile(reportName string, billingFile string) error {
	var end = false

	l.reportName = reportName
	l.report = l.openBillingReport(billingFile)
	defer l.reportFd.Close()

	for end == false {
		record := &billingRecord{}
		values, readErr := l.report.csvReader.Read()
		if readErr == io.EOF {
			end = true
		} else if readErr != nil {
			fmt.Printf("error while reading %s: %v", l.reportName, readErr)
			continue
		}
		if parseErr := l.parseRecord(values, record, l.report); parseErr != nil {
			continue
		}
		if saveErr := l.saveRecord(record, end); saveErr != nil {
			return saveErr
		}
	}
	return nil
}

type billingLoader struct {
	reportName string
	report     *billingReport
	reportFd   *os.File
	json       []byte
	injector   Injector
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
		if i < 0 || i > len(in) {
			continue
		}
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
	return nil
}

func (l *billingLoader) assignField(record *billingRecord, field string, value interface{}, tag bool) error {
	if tag == false {
		recordField := reflect.ValueOf(record).Elem().FieldByName(field)
		if !recordField.IsValid() {
			return fmt.Errorf("No such field '%s'", field)
		}
		if !recordField.CanSet() {
			return fmt.Errorf("Field '%s' cannot be set.", field)
		}
		valueTyped := reflect.ValueOf(value)
		if recordField.Type() != valueTyped.Type() {
			return fmt.Errorf("Field '%s' didn't have type %s.", field, valueTyped.Type().Name())
		}
		recordField.Set(valueTyped)
	} else {
		for k, v := range value.(map[string]string) {
			record.Tags[k] = v
		}
	}
	return nil
}

func (l *billingLoader) openBillingReport(billingFile string) *billingReport {
	var dateMonthMatcher = regexp.MustCompile(`\d{4}-\d{2}`)

	invoicePeriod, err := time.Parse("2006-02", dateMonthMatcher.FindString(billingFile))
	if err != nil {
		panic(err.Error())
	}
	file, err := os.Open(billingFile)
	if err != nil {
		panic(err.Error())
	}
	l.reportFd = file
	reader := csv.NewReader(file)
	report := &billingReport{
		CsvFileName:   billingFile,
		InvoicePeriod: invoicePeriod,
		csvReader:     reader}

	fields, err := reader.Read()
	if err != nil {
		panic(err.Error())
	}
	report.Fields = fields
	return report
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

func (l *billingLoader) saveRecord(in chan *billingRecord) {
	for record := range in {
		record.ReportName = l.reportName
		record.Id = uuid.NewV4().String()
		val, err := json.MarshalIndent(*record, "", "  ")
		if err == nil {
			fmt.Printf("%s ===================\n", string(val))
		}
		//		err = l.injector.CreateRecord(record)
		//if err != nil {
		//	fmt.Printf("Error during import: %v\n", err)
		//}
	}
}
