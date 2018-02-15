package raws_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/cycloidio/raws"
)

func ExampleAWSReader_GetLoadBalancersV2() {
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

	var dli = &elbv2.DescribeLoadBalancersInput{}
	var _, awsErrs = awsr.GetLoadBalancersV2(ctx, dli)
	if awsErrs != nil {
		fmt.Printf("Error: %+v\n", awsErrs)
		return
	}
	// Output:
}
