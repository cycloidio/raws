package raws_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/cycloidio/raws"
)

func ExampleAWSReader_GetInstances() {
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

	var dii = &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{Name: aws.String("tag:env"), Values: []*string{aws.String("prod")}},
		},
	}

	var _, awsErrs = awsr.GetInstances(ctx, dii)
	if awsErrs != nil {
		fmt.Printf("Error: %+v\n", awsErrs)
		return
	}

	// Output:
}

func ExampleAWSReader_GetVpcs() {
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

	var dvi = &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{Name: aws.String("state"), Values: []*string{aws.String("available")}},
		},
	}

	var _, awsErrs = awsr.GetVpcs(ctx, dvi)
	if awsErrs != nil {
		fmt.Printf("Error: %+v\n", awsErrs)
		return
	}

	// Output:
}

func ExampleAWSReader_GetSubnets() {
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

	var dsi = &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("available-ip-address-count"),
				Values: []*string{aws.String("251")},
			},
		},
	}

	var _, awsErrs = awsr.GetSubnets(ctx, dsi)
	if awsErrs != nil {
		fmt.Printf("Error: %+v\n", awsErrs)
		return
	}
	// Output:
}

func ExampleAWSReader_GetSecurityGroups() {
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

	var dsi = &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("group-name"),
				Values: []*string{aws.String("test_allow_metrics")},
			},
		},
	}

	var _, awsErrs = awsr.GetSecurityGroups(ctx, dsi)
	if awsErrs != nil {
		fmt.Printf("Error: %+v\n", awsErrs)
		return
	}
	// Output:
}

func ExampleAWSReader_GetImages() {
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

	var dsi = &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("architecture"),
				Values: []*string{aws.String("x86_64")},
			},
		},
	}

	var _, awsErrs = awsr.GetImages(ctx, dsi)
	if awsErrs != nil {
		fmt.Printf("Error: %+v\n", awsErrs)
		return
	}
	// Output:
}

func ExampleAWSReader_GetVolumes() {
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

	var dvi = &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("size"),
				Values: []*string{aws.String("10")},
			},
		},
	}

	var _, awsErrs = awsr.GetVolumes(ctx, dvi)
	if awsErrs != nil {
		fmt.Printf("Error: %+v\n", awsErrs)
		return
	}
	// Output:
}

func ExampleAWSReader_GetSnapshots() {
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

	var dsi = &ec2.DescribeSnapshotsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("progress"),
				Values: []*string{aws.String("100%")},
			},
		},
	}

	var _, awsErrs = awsr.GetSnapshots(ctx, dsi)
	if awsErrs != nil {
		fmt.Printf("Error: %+v\n", awsErrs)
		return
	}
	// Output:
}
