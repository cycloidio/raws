package billing

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func TestNewInjector(t *testing.T) {
	var mockedDyn dynamodbiface.DynamoDBAPI = mockDynamodb{}

	i := &billingInjector{
		dynamoSvc:   mockedDyn,
		requests:    make(map[string][]*dynamodb.WriteRequest),
		errorsCount: 0,
	}
	ci := NewInjector(mockedDyn)
	if !reflect.DeepEqual(i, ci) {
		t.Errorf("NewInjector: received=%+v | expected=%+v",
			ci, i)
	}
}

func TestBillingInjector_CreateReport(t *testing.T) {
	tests := []struct {
		name          string
		mockedDyn     dynamodbiface.DynamoDBAPI
		filename      string
		hash          string
		expectedError error
	}{
		{name: "no error while creating report",
			mockedDyn: mockDynamodb{
				pio: &dynamodb.PutItemOutput{},
				pie: nil,
			},
			filename:      "test",
			hash:          "test-hash",
			expectedError: nil,
		},
		{name: "error while creating report",
			mockedDyn: mockDynamodb{
				pio: &dynamodb.PutItemOutput{},
				pie: errors.New("cannot create"),
			},
			filename:      "test",
			hash:          "test-hash",
			expectedError: NewDynamoDBError(errors.New("cannot create")),
		},
	}

	for j, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &billingInjector{
				dynamoSvc: tt.mockedDyn,
			}
			err := i.CreateReport(tt.filename, tt.hash)
			checkErrors(t, tt.name, j, err, tt.expectedError)
		})
	}
}

func TestBillingInjector_CreateRecord(t *testing.T) {
	tests := []struct {
		name          string
		record        *billingRecord
		mockedDyn     dynamodbiface.DynamoDBAPI
		expectedError error
	}{
		{name: "no error while creating record",
			mockedDyn: mockDynamodb{
				pio: &dynamodb.PutItemOutput{},
				pie: nil,
			},
			record:        &billingRecord{},
			expectedError: nil,
		},
		{name: "error while creating record",
			mockedDyn: mockDynamodb{
				pio: &dynamodb.PutItemOutput{},
				pie: errors.New("cannot create"),
			},
			record:        &billingRecord{},
			expectedError: NewDynamoDBError(errors.New("cannot create")),
		},
	}

	for j, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &billingInjector{
				dynamoSvc: tt.mockedDyn,
			}
			err := i.CreateRecord(tt.record)
			checkErrors(t, tt.name, j, err, tt.expectedError)
		})
	}
}

func TestBillingInjector_MaxRecords(t *testing.T) {
	i := billingInjector{}
	max := i.MaxRecords()
	if max != maxRequestSize {
		t.Errorf("MaxRecords doesn't match: received=%d | expected=%d",
			max, maxRequestSize)
	}
}

func TestBillingInjector_createRecordIdList(t *testing.T) {
	var expectedList = []string{"1", "2", "3"}

	i := billingInjector{}
	rs := []*billingRecord{{RecordId: "0"}, {RecordId: "1"}, {RecordId: "2"}, {RecordId: ""}, {RecordId: "3"}}
	list := i.createRecordIdList(rs)
	if !reflect.DeepEqual(list, expectedList) {
		t.Errorf("createRecordIdList: received=%+v | expected=%+v",
			list, expectedList)
	}
}

func TestBillingInjector_createRecordIdListFromRequest(t *testing.T) {
	var expectedList = []string{"1", "2", "3"}

	i := billingInjector{
		requests: make(map[string][]*dynamodb.WriteRequest),
	}
	rs := []*billingRecord{{RecordId: "0"}, {RecordId: "1"}, {RecordId: "2"}, {RecordId: ""}, {RecordId: "3"}}
	err := i.createRequest(rs)
	if err != nil {
		t.Errorf("createRecordIdListFromRequest - issue while creating request: received=%+v", err)
	}

	list := i.createRecordIdListFromRequest()
	if !reflect.DeepEqual(list, expectedList) {
		t.Errorf("createRecordIdListFromRequest: received=%+v | expected=%+v",
			list, expectedList)
	}
}

func TestBillingInjector_CreateRecords(t *testing.T) {
	var sampleRecord *billingRecord = &billingRecord{RecordId: "1"}
	av, err := dynamodbattribute.MarshalMap(sampleRecord)
	if err != nil {
		t.Errorf("CreateRecords: couldn't create sample received=%+v", err)
	}
	sampleWriteRequest := &dynamodb.WriteRequest{
		PutRequest: &dynamodb.PutRequest{
			Item: av,
		},
	}

	tests := []struct {
		name              string
		records           []*billingRecord
		mockedDyn         dynamodbiface.DynamoDBAPI
		expectedList      []string
		expectedProcessed int
		expectedError     error
	}{
		{name: "no error while creating records",
			mockedDyn: mockDynamodb{
				bwio: &dynamodb.BatchWriteItemOutput{},
				bwie: nil,
			},
			records:           []*billingRecord{{RecordId: "1"}},
			expectedError:     nil,
			expectedProcessed: 1,
			expectedList:      nil,
		},
		{name: "error while creating records",
			mockedDyn: mockDynamodb{
				bwio: &dynamodb.BatchWriteItemOutput{
					UnprocessedItems: map[string][]*dynamodb.WriteRequest{
						billingRecordTableName: {sampleWriteRequest},
					},
				},
				bwie: errors.New("cannot create"),
			},
			records:           []*billingRecord{{RecordId: "1"}},
			expectedError:     NewDynamoDBError(errors.New("cannot create")),
			expectedProcessed: 0,
			expectedList:      []string{"1"},
		},
	}

	for j, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &billingInjector{
				dynamoSvc: tt.mockedDyn,
				requests:  make(map[string][]*dynamodb.WriteRequest),
			}
			list, processed, err := i.CreateRecords(tt.records)
			checkErrors(t, tt.name, j, err, tt.expectedError)
			if processed != tt.expectedProcessed {
				t.Errorf("CreateRecords: received=%d | expected=%d",
					processed, tt.expectedProcessed)
			}
			if !reflect.DeepEqual(list, tt.expectedList) {
				t.Errorf("CreateRecords recordIds invalid: received=%+v | expected=%+v",
					list, tt.expectedList)
			}
		})
	}
}
