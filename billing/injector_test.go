package billing

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func TestNewInjector(t *testing.T) {
	var mockedDyn dynamodbiface.DynamoDBAPI = mockDynamodb{}

	i := &billingInjector{
		dynamoSvc: mockedDyn,
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
			expectedError: errors.New("cannot create"),
		},
	}

	for j, tt := range tests {
		i := &billingInjector{
			dynamoSvc: tt.mockedDyn,
		}
		err := i.CreateReport(tt.filename, tt.hash)
		checkErrors(t, tt.name, j, err, tt.expectedError)
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
			expectedError: errors.New("cannot create"),
		},
	}

	for j, tt := range tests {
		i := &billingInjector{
			dynamoSvc: tt.mockedDyn,
		}
		err := i.CreateRecord(tt.record)
		checkErrors(t, tt.name, j, err, tt.expectedError)
	}
}
