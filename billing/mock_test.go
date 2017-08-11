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
}

func (m mockDownloader) Download(bucket string, filename string, dest string) (string, error) {
	return m.ds, m.de
}

type mockLoader struct {
	Loader
}

func (m mockLoader) ProcessFile(reportName string, billingFile string) {
	return
}

type mockInjector struct {
	Injector

	// CreateRecord returned element
	crde error

	// CreateReport returned element
	crpe error
}

func (m mockInjector) CreateRecord(record *billingRecord) error {
	return m.crde
}
func (m mockInjector) CreateReport(filename string, hash string) error {
	return m.crpe
}

// mockAWSReader used for testing purposes
type mockAWSReader struct {
	raws.AWSReader

	// GetAccountID returned element
	accountID string

	// ListObjects returned elements
	loo []*s3.ListObjectsOutput
	loe raws.Errs

	// DownloadObject returned elements
	doi int64
	doe error
}

func (m mockAWSReader) GetAccountID() string {
	return m.accountID
}

func (m mockAWSReader) ListObjects(input *s3.ListObjectsInput) ([]*s3.ListObjectsOutput, raws.Errs) {
	return m.loo, m.loe
}

func (m mockAWSReader) DownloadObject(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error) {
	return m.doi, m.doe
}

func checkErrors(t *testing.T, name string, index int, err error, expected error) {
	if err != nil && !reflect.DeepEqual(err, expected) {
		t.Errorf("%s [%d] - errors: received=%+v | expected=%+v",
			name, index, err.Error(), expected.Error())
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
}

func (m mockDynamodb) GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return m.geo, m.gee
}

func (m mockDynamodb) PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return m.pio, m.pie
}
