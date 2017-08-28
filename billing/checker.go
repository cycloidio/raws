package billing

import (
	"fmt"

	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cycloidio/raws"
)

type Checker interface {
	Check(bucket string, filename string) (bool, error)
	AlreadyPresent() (bool, string)
}

type billingChecker struct {
	s3Connector raws.AWSReader
	dynSvc      dynamodbiface.DynamoDBAPI
	oldMd5      string
	newMd5      string
}

func NewChecker(s3connector raws.AWSReader, dynamoDB dynamodbiface.DynamoDBAPI) Checker {
	return &billingChecker{
		s3Connector: s3connector,
		dynSvc:      dynamoDB,
		oldMd5:      "",
		newMd5:      "",
	}
}

func (c *billingChecker) Check(bucket string, filename string) (bool, error) {
	err := c.getDynamoEntry(filename)
	if err != nil {
		return false, err
	}
	err = c.getS3Entry(bucket, filename)
	if err != nil {
		return false, err
	}
	if c.newMd5 == c.oldMd5 {
		return false, nil
	}
	return true, nil
}

func (c *billingChecker) AlreadyPresent() (bool, string) {
	if c.oldMd5 == c.newMd5 {
		return true, c.newMd5
	}
	return false, c.newMd5
}

func (c *billingChecker) getS3Entry(bucket string, filename string) error {
	inputs := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(filename),
	}

	objectsOutput, err := c.s3Connector.ListObjects(inputs)
	if err != nil {
		return err
	}
	if len(objectsOutput) != 1 {
		return fmt.Errorf("found too many objects matching (%d)", len(objectsOutput))
	}
	if objectsOutput[0].Contents == nil || len(objectsOutput[0].Contents) == 0 {
		return errors.New("s3 entry doesn't have 'Contents' attribute")
	}
	etag := *objectsOutput[0].Contents[0].ETag
	c.newMd5 = etag[1 : len(etag)-1]
	return nil
}

func (c *billingChecker) getDynamoEntry(filename string) error {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			billingReportNameField: {
				S: aws.String(filename),
			},
		},
		TableName: aws.String(billingReportTableName),
	}
	result, err := c.dynSvc.GetItem(input)
	if err != nil {
		return err
	}
	if result == nil || len(result.Item) == 0 {
		return nil
	}
	if val, ok := result.Item[billingReportMd5Field]; ok {
		c.oldMd5 = *val.S
		return nil
	}
	return fmt.Errorf("no '%s' field present for the entity", billingReportMd5Field)
}
