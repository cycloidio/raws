package billing

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Injector interface {
	CreateRecord(record *billingRecord) error
	CreateReport(filename string, hash string) error
}

type billingInjector struct {
	dynamoSvc *dynamodb.DynamoDB
}

func NewInjector(dynamoSvc *dynamodb.DynamoDB) Injector {
	return &billingInjector{
		dynamoSvc: dynamoSvc,
	}
}

func (i *billingInjector) CreateRecord(record *billingRecord) error {
	av, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		return err
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
	fmt.Printf("Item: %v\n", input)
	return err
}
