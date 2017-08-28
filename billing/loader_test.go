package billing

import (
	"errors"
	"reflect"
	"testing"

	"time"

	"strconv"

	"encoding/csv"

	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func TestBillingLoader_NewLoader(t *testing.T) {
	var mockedDyn dynamodbiface.DynamoDBAPI = mockDynamodb{}
	var i Injector

	i = NewInjector(mockedDyn)
	b := &billingLoader{
		json:     []byte{},
		injector: i,
		result:   newStats(),
		reportFd: nil,
	}

	cb := NewLoader(i)
	if !reflect.DeepEqual(b, cb) {
		t.Errorf("NewLoader: received=%+v | expected=%+v",
			cb, b)
	}
}

func TestBillingLoader_ProcessFile(t *testing.T) {
	tests := []struct {
		name            string
		reportName      string
		filePath        string
		mockedInjector  mockInjector
		expectedError   error
		expectedRecords []string
	}{
		{
			name:       "file is valid but not name",
			reportName: "invalid-filename",
			filePath:   "test/csvs/invalid-name.csv",
			mockedInjector: mockInjector{
				crss: nil,
				crsi: 10,
				crse: nil,
				mri:  10,
			},
			expectedError: &time.ParseError{
				Value:      "",
				Layout:     "2006-02",
				ValueElem:  "",
				LayoutElem: "2006",
			},
			expectedRecords: nil,
		},
		{
			name:       "file has only good values",
			reportName: "valid-filenamem",
			filePath:   "test/csvs/valid-2017-07.csv",
			mockedInjector: mockInjector{
				crss: nil,
				crsi: 10,
				crse: nil,
				mri:  10,
			},
			expectedError:   nil,
			expectedRecords: nil,
		},
		{
			name:       "file has mixed errors/good values",
			reportName: "valid-filename",
			filePath:   "test/csvs/mixed-errors-2017-07.csv",
			mockedInjector: mockInjector{
				crss: nil,
				crsi: 10,
				crse: nil,
				mri:  10,
			},
			expectedError: NewConvertError(&strconv.NumError{
				Func: "ParseInt",
				Num:  "",
				Err:  errors.New("invalid syntax"),
			}),
			expectedRecords: nil,
		},
		{
			name:       "file is valid but not name",
			reportName: "valid-filename",
			filePath:   "test/csvs/only-errors-2017-07.csv",
			mockedInjector: mockInjector{
				crss: nil,
				crsi: 10,
				crse: nil,
				mri:  10,
			},
			expectedError: NewConvertError(&strconv.NumError{
				Func: "ParseInt",
				Num:  "",
				Err:  errors.New("invalid syntax"),
			}),
			expectedRecords: nil,
		},
		{
			name:       "file is valid but failed to inject",
			reportName: "valid-filename",
			filePath:   "test/csvs/valid-2017-07.csv",
			mockedInjector: mockInjector{
				crss: nil,
				crsi: 10,
				crse: errors.New("failed to inject"),
				mri:  10,
			},
			expectedError:   errors.New("failed to inject"),
			expectedRecords: nil,
		},
		{
			name:       "file has invalid column number",
			reportName: "invalid-file",
			filePath:   "test/csvs/invalid-column-2017-07.csv",
			mockedInjector: mockInjector{
				crss: nil,
				crsi: 10,
				crse: nil,
				mri:  10,
			},
			expectedError: NewCSVError(&csv.ParseError{
				Line:   2,
				Column: 0,
				Err:    errors.New("wrong number of fields in line"),
			}),
			expectedRecords: nil,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &billingLoader{
				injector: tt.mockedInjector,
				result:   newStats(),
			}
			records, err := m.ProcessFile(tt.reportName, tt.filePath)
			checkErrors(t, tt.name, i, err, tt.expectedError)
			if !reflect.DeepEqual(records, tt.expectedRecords) {
				t.Errorf("Invalid records: received=%+v | expected=%+v",
					records, tt.expectedRecords)
			}
		})
	}

}

func TestBillingLoader_parseRecord(t *testing.T) {
	t.Run("parse valid data", func(t *testing.T) {
		var fields = []string{
			"InvoiceID", "PayerAccountId", "LinkedAccountId", "RecordType",
			"RecordId", "ProductName", "RateId", "SubscriptionId", "PricingPlanId",
			"UsageType", "Operation", "AvailabilityZone", "ReservedInstance",
			"ItemDescription", "UsageStartDate", "UsageEndDate", "UsageQuantity",
			"BlendedRate", "BlendedCost", "UnBlendedRate", "UnBlendedCost", "ResourceId",
			"user:env", "user:project",
		}
		var inputs = []string{
			"52511536", "661913936052", "137871057171", "LineItem",
			"35330993584143683082238491", "Amazon CloudFront", "3605351",
			"130980400", "509640", "AP-DataTransfer-Out-Bytes", "GET", "",
			"N", "$0.000 per GB - data transfer out under the global monthly free tier",
			"2015-04-30 14:00:00", "2015-04-30 15:00:00", "0.01101368", "0.0000000000",
			"0.00000000", "0.0000000000", "0.00000000", "E29NSV4ST80EBC", "prod", "test",
		}
		var expectedError error
		var expectedRecord = &billingRecord{
			InvoiceID:        "52511536",
			PayerAccountId:   661913936052,
			LinkedAccountId:  137871057171,
			RecordType:       "LineItem",
			RecordId:         "35330993584143683082238491",
			ProductName:      "Amazon CloudFront",
			RateId:           3605351,
			SubscriptionId:   130980400,
			PricingPlanId:    509640,
			UsageType:        "AP-DataTransfer-Out-Bytes",
			Operation:        "GET",
			AvailabilityZone: "",
			ReservedInstance: "N",
			ItemDescription:  "$0.000 per GB - data transfer out under the global monthly free tier",
			UsageStartDate:   "2015-04-30T14:00:00Z",
			UsageEndDate:     "2015-04-30T15:00:00Z",
			UsageQuantity:    0.01101368,
			BlendedRate:      0.0000000000,
			BlendedCost:      0.00000000,
			UnBlendedRate:    0.0000000000,
			UnBlendedCost:    0.00000000,
			ResourceId:       "E29NSV4ST80EBC",
			Tags: map[string]string{
				"user_env":     "prod",
				"user_project": "test",
			},
		}
		record := &billingRecord{}
		report := &billingReport{
			Fields: fields,
		}
		m := &billingLoader{
			result: newStats(),
		}
		err := m.parseRecord(inputs, record, report)
		checkErrors(t, "create-valid-record", 0, err, expectedError)
		if !reflect.DeepEqual(record, expectedRecord) {
			t.Errorf("%s - record differ: received=%+v | expected=%+v",
				"create-valid-record", record, expectedRecord)
		}
	})
	t.Run("parse invalid data", func(t *testing.T) {
		var fields = []string{
			"InvoiceID", "PayerAccountId", "LinkedAccountId", "RecordType",
			"RecordId", "ProductName", "RateId", "SubscriptionId", "PricingPlanId",
			"UsageType", "Operation", "AvailabilityZone", "ReservedInstance",
			"ItemDescription", "UsageStartDate", "UsageEndDate", "UsageQuantity",
			"BlendedRate", "BlendedCost", "UnBlendedRate", "UnBlendedCost", "ResourceId",
			"user:env", "user:project",
		}
		var inputs = []string{
			"52511536", "661913936052", "137871057171", "LineItem",
			"", "Amazon CloudFront", "3605351",
			"130980400", "509640", "AP-DataTransfer-Out-Bytes", "GET", "",
			"N", "$0.000 per GB - data transfer out under the global monthly free tier",
			"2015-04-30 14:00:00", "2015-04-30 15:00:00", "0.01101368", "0.0000000000",
			"0.00000000", "0.0000000000", "0.00000000", "E29NSV4ST80EBC", "prod", "test",
		}
		var expectedError = fmt.Errorf("no recordId found for this entry: %v", inputs)
		var expectedRecord = &billingRecord{
			InvoiceID:        "52511536",
			PayerAccountId:   661913936052,
			LinkedAccountId:  137871057171,
			RecordType:       "LineItem",
			RecordId:         "",
			ProductName:      "Amazon CloudFront",
			RateId:           3605351,
			SubscriptionId:   130980400,
			PricingPlanId:    509640,
			UsageType:        "AP-DataTransfer-Out-Bytes",
			Operation:        "GET",
			AvailabilityZone: "",
			ReservedInstance: "N",
			ItemDescription:  "$0.000 per GB - data transfer out under the global monthly free tier",
			UsageStartDate:   "2015-04-30T14:00:00Z",
			UsageEndDate:     "2015-04-30T15:00:00Z",
			UsageQuantity:    0.01101368,
			BlendedRate:      0.0000000000,
			BlendedCost:      0.00000000,
			UnBlendedRate:    0.0000000000,
			UnBlendedCost:    0.00000000,
			ResourceId:       "E29NSV4ST80EBC",
			Tags: map[string]string{
				"user_env":     "prod",
				"user_project": "test",
			},
		}
		record := &billingRecord{}
		report := &billingReport{
			Fields: fields,
		}
		m := &billingLoader{
			result: newStats(),
		}
		err := m.parseRecord(inputs, record, report)
		checkErrors(t, "create-valid-record", 0, err, expectedError)
		if !reflect.DeepEqual(record, expectedRecord) {
			t.Errorf("%s - record differ: received=%+v | expected=%+v",
				"create-valid-record", record, expectedRecord)
		}
	})
}
