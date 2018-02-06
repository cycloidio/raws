package raws_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cycloidio/raws"
)

func ExampleAWSReader_ListBuckets() {
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

	var lbi = &s3.ListBucketsInput{}
	var _, awsErrs = awsr.ListBuckets(ctx, lbi)
	if awsErrs != nil {
		fmt.Printf("Error: %+v\n", awsErrs)
		return
	}
	// Output:
}

func ExampleAWSReader_GetBucketTags() {
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

	var lbi = &s3.ListBucketsInput{}
	var buckets, awsErrs = awsr.ListBuckets(ctx, lbi)
	if awsErrs != nil {
		fmt.Printf("Error: %+v\n", awsErrs)
		return
	}

	var bucketName = ""
	if len(buckets) > 0 {
		for _, s := range buckets {
			if len(s.Buckets) > 0 {
				bucketName = *s.Buckets[0].Name
				break
			}
		}

		if bucketName == "" {
			return
		}
	}

	fmt.Println(bucketName)

	var gbti = &s3.GetBucketTaggingInput{
		Bucket: aws.String(bucketName),
	}
	var tags map[string]s3.GetBucketTaggingOutput
	tags, awsErrs = awsr.GetBucketTags(ctx, gbti)
	if awsErrs != nil {
		fmt.Printf("Partial results: %+v\n", tags)
		fmt.Printf("Error: %+v\n", awsErrs)
		return
	}
	// Output:
}
