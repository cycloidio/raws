package billing

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/twinj/uuid"
)

type Loader struct {
	reportName  string
	report      *billingReport
	concurrency int
	wg          sync.WaitGroup
	json        []byte
	injector    Injector
}

func NewLoader(injector Injector, concurrency int) *Loader {
	return &Loader{
		concurrency: concurrency,
		wg:          sync.WaitGroup{},
		json:        []byte{},
		injector:    injector,
	}
}

func (l *Loader) ProcessFile(reportName string, billingFile string) {
	l.reportName = reportName
	l.report = l.openBillingReport(billingFile)

	valuesSink := make(chan []string)
	reportSinks := make([]chan *billingRecord, 0, 0)

	for i := 0; i < runtime.GOMAXPROCS(l.concurrency); i++ {
		reportSink := make(chan *billingRecord)
		reportSinks = append(reportSinks, reportSink)
		l.wg.Add(2)
		go l.parseRecord(valuesSink, reportSink, l.report)
		go l.saveRecord(reportSink)
	}

	for {
		values, err := l.report.csvReader.Read()
		if err == io.EOF {
			break
		}
		valuesSink <- values
	}
	close(valuesSink)
	l.wg.Wait()
}

type fieldMapper map[string]map[interface{}]map[string]interface{}

type billingReport struct {
	CsvFileName   string
	InvoicePeriod time.Time
	ReportType    int
	Fields        []string
	csvReader     *csv.Reader
	Mapper        fieldMapper
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

func (l *Loader) parseRecord(in chan []string, out chan *billingRecord, report *billingReport) {
	var tagMatcher = regexp.MustCompile("(user|aws):.*")
	var fieldTypes = map[string]func(s string, report *billingReport) interface{}{
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

	for values := range in {
		var record = &billingRecord{}
		record.Tags = make(map[string]string)
		for i, field := range report.Fields {
			if f, ok := fieldTypes[field]; ok {
				l.assignField(record, field, f(values[i], report), false)
			} else if tagMatcher.MatchString(field) {
				value := map[string]string{l.parseTag(field, report).(string): values[i]}
				l.assignField(record, field, value, true)
			} else {
				l.assignField(record, field, values[i], false)
			}
		}
		out <- record
	}
	close(out)
	l.wg.Done()
}

func (l *Loader) assignField(record *billingRecord, field string, value interface{}, tag bool) error {
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

func (l *Loader) openBillingReport(billingFile string) *billingReport {
	var dateMonthMatcher = regexp.MustCompile(`\d{4}-\d{2}`)

	invoicePeriod, err := time.Parse("2006-02", dateMonthMatcher.FindString(billingFile))
	if err != nil {
		panic(err.Error())
	}
	file, err := os.Open(billingFile)
	if err != nil {
		panic(err.Error())
	}
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

func (l *Loader) parseTag(s string, report *billingReport) interface{} {
	tagParts := strings.Split(s, ":")
	return strings.Join(tagParts, "_")
}

func (l *Loader) parseInt(s string, report *billingReport) interface{} {
	value, _ := strconv.ParseInt(s, 0, 0)
	return value
}

func (l *Loader) parseUint(s string, report *billingReport) interface{} {
	value, _ := strconv.ParseUint(s, 0, 0)
	return value
}

func (l *Loader) parseFloat(s string, report *billingReport) interface{} {
	value, _ := strconv.ParseFloat(s, 0)
	return value
}

func (l *Loader) parseDate(s string, report *billingReport) interface{} {
	var returnTime time.Time
	var err error
	switch s {
	case "":
		returnTime = report.InvoicePeriod
	default:
		returnTime, err = time.Parse("2006-01-02 15:04:05", s)
		if err != nil {
			panic(err.Error())
		}
	}
	return returnTime.Format(time.RFC3339)
}

func (l *Loader) saveRecord(in chan *billingRecord) {
	for record := range in {
		record.ReportName = l.reportName
		record.Id = uuid.NewV4().String()
		val, err := json.MarshalIndent(*record, "", "  ")
		if err == nil {
			fmt.Printf("%s ===================\n", string(val))
		}
		err = l.injector.CreateRecord(record)
		if err != nil {
			fmt.Printf("Error during import: %v\n", err)
		}
	}
	l.wg.Done()
}
