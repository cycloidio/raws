package billing

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type Injector interface {
	CreateRecord(record *billingRecord) error
	CreateReport(filename string, hash string) error
}

type billingInjector struct {
	dynamoSvc dynamodbiface.DynamoDBAPI
}

func NewInjector(dynamoSvc dynamodbiface.DynamoDBAPI) Injector {
	return &billingInjector{
		dynamoSvc: dynamoSvc,
	}
}

func (i *billingInjector) CreateRecord(record *billingRecord) error {
	av, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		return err
	}

	if record.RecordId == "0" || record.RecordId == "" {
		return nil
	}
	_, err = i.dynamoSvc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(billingRecordTableName),
		Item:      av,
	})
	return err
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
		},
		TableName:              aws.String(billingReportTableName),
		ReturnConsumedCapacity: aws.String("TOTAL"),
	}
	_, err := i.dynamoSvc.PutItem(input)
	return err
}
