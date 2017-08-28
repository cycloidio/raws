package billing

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type Injector interface {
	CreateRecord(record *billingRecord) error
	CreateRecords(records []*billingRecord) ([]string, int, error)
	CreateReport(filename string, hash string) error
	MaxRecords() int
}

const (
	maxRequestSize int = 25

	billingReportTableName   = "billing-reports"
	billingReportNameField   = "name"
	billingReportMd5Field    = "md5"
	billingReportErrorsField = "errors"

	billingRecordTableName = "billing-records"
)

type billingInjector struct {
	dynamoSvc   dynamodbiface.DynamoDBAPI
	requests    map[string][]*dynamodb.WriteRequest
	errorsCount int
}

func NewInjector(dynamoSvc dynamodbiface.DynamoDBAPI) Injector {
	return &billingInjector{
		dynamoSvc:   dynamoSvc,
		requests:    make(map[string][]*dynamodb.WriteRequest),
		errorsCount: 0,
	}
}

func (i *billingInjector) MaxRecords() int {
	return maxRequestSize
}

func (i *billingInjector) CreateRecords(records []*billingRecord) ([]string, int, error) {
	var dynErr error
	var result *dynamodb.BatchWriteItemOutput
	var initial int

	err := i.createRequest(records)
	if err != nil {
		i.errorsCount += len(records)
		return i.createRecordIdList(records), initial, err
	}
	initial = len(i.requests[billingRecordTableName])
	for j := 0; j < 3; j++ {
		result, dynErr = i.dynamoSvc.BatchWriteItem(&dynamodb.BatchWriteItemInput{
			RequestItems: i.requests,
		})
		if dynErr == nil {
			break
		}
		i.requests = result.UnprocessedItems
	}
	if dynErr != nil {
		i.errorsCount += len(i.requests[billingRecordTableName])
		return i.createRecordIdListFromRequest(),
			initial - len(i.requests[billingRecordTableName]),
			NewDynamoDBError(dynErr)
	}
	return nil, initial, nil
}

func (i *billingInjector) CreateRecord(record *billingRecord) error {
	av, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		i.errorsCount++
		return NewDynamoDBError(err)
	}

	_, err = i.dynamoSvc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(billingRecordTableName),
		Item:      av,
	})
	if err != nil {
		i.errorsCount++
		return NewDynamoDBError(err)
	}
	return nil
}

func (i *billingInjector) CreateReport(filename string, hash string) error {
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			billingReportNameField: {
				S: aws.String(filename),
			},
			billingReportMd5Field: {
				S: aws.String(hash),
			},
			billingReportErrorsField: {
				N: aws.String(strconv.Itoa(i.errorsCount)),
			},
		},
		TableName:              aws.String(billingReportTableName),
		ReturnConsumedCapacity: aws.String("TOTAL"),
	}
	_, err := i.dynamoSvc.PutItem(input)
	if err != nil {
		return NewDynamoDBError(err)
	}
	return nil
}

func (i *billingInjector) createRecordIdList(records []*billingRecord) []string {
	var recordIds []string

	for _, record := range records {
		if record.RecordId == "0" || record.RecordId == "" {
			continue
		}
		recordIds = append(recordIds, record.RecordId)
	}
	return recordIds
}

func (i *billingInjector) createRecordIdListFromRequest() []string {
	var recordIds []string

	for _, req := range i.requests[billingRecordTableName] {
		if req.PutRequest != nil || req.PutRequest.Item != nil {
			if val, ok := req.PutRequest.Item["RecordId"]; ok {
				if val == nil || val.S == nil || *val.S == "0" || *val.S == "" {
					continue
				}
				recordId := *val.S
				recordIds = append(recordIds, recordId)
			}
		}
	}
	return recordIds
}

func (i *billingInjector) createRequest(records []*billingRecord) error {
	i.requests[billingRecordTableName] = nil

	for _, record := range records {
		av, err := dynamodbattribute.MarshalMap(record)

		if err != nil {
			i.requests[billingRecordTableName] = nil
			return NewDynamoDBError(err)
		}

		request := &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: av,
			},
		}
		i.requests[billingRecordTableName] = append(i.requests[billingRecordTableName], request)
	}

	return nil
}
