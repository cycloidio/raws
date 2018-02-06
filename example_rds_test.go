package raws_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/cycloidio/raws"
)

func ExampleAWSReader_GetDBInstances() {
	var accessKey, secretKey, err = getAWSKeys()
	if err != nil {
		// When no keys the example doesn't run, but if the env vars are set
		// then it runs and make the AWS SDK calls.
		return
	}

	var ctx = context.Background()
	var awsr raws.AWSReader
	awsr, err = raws.NewAWSReader(ctx, accessKey, secretKey, []string{"eu-*"}, nil)
	if err != nil {
		fmt.Printf("Error while getting NewConnector: %s\n", err.Error())
		return
	}

	var ddbi = &rds.DescribeDBInstancesInput{}
	var _, awsErrs = awsr.GetDBInstances(ctx, ddbi)
	if awsErrs != nil {
		fmt.Printf("Error: %+v\n", awsErrs)
		return
	}
	// Output:
}
