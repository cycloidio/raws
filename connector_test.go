package raws

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sts"
)

func TestConnector_setRegion(t *testing.T) {
	var ec2Regions = &ec2.DescribeRegionsOutput{
		Regions: []*ec2.Region{
			{RegionName: aws.String("eu-west-1")},
			{RegionName: aws.String("eu-west-2")},
			{RegionName: aws.String("eu-central-1")},
			{RegionName: aws.String("us-west-1")},
		},
	}

	tests := []struct {
		name            string
		mocked          mockEC2
		regionsInput    []string
		expectedRegions []string
		expectedError   error
	}{
		{name: "no regions given",
			mocked: mockEC2{
				dro:   ec2Regions,
				drerr: nil,
			},
			regionsInput:    []string{},
			expectedRegions: nil,
			expectedError:   errors.New("at least one region name is required"),
		},
		{name: "match all eu regions",
			mocked: mockEC2{
				dro:   ec2Regions,
				drerr: nil,
			},
			regionsInput:    []string{"eu-*"},
			expectedRegions: []string{"eu-west-1", "eu-west-2", "eu-central-1"},
			expectedError:   nil,
		},
		{name: "match nothing",
			mocked: mockEC2{
				dro:   ec2Regions,
				drerr: nil,
			},
			regionsInput:    []string{"matchnothing-*"},
			expectedRegions: nil,
			expectedError:   errors.New("found 0 regions matching: [matchnothing-*]"),
		},
		{name: "error in describe regions",
			mocked: mockEC2{
				dro:   nil,
				drerr: errors.New("fail"),
			},
			regionsInput:    []string{"matchnothing-*"},
			expectedRegions: nil,
			expectedError:   errors.New("fail"),
		}}

	var ctx = context.Background()

	for i, tt := range tests {
		c := &connector{}
		err := c.setRegions(ctx, tt.mocked, tt.regionsInput)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(c.regions, tt.expectedRegions) {
			t.Errorf("%s [%d] - regions: received=%+v | expected=%+v",
				tt.name, i, c.regions, tt.expectedRegions)
		}
	}
}

func TestConnector_setAccountID(t *testing.T) {
	tests := []struct {
		name          string
		mocked        mockSTS
		expectedID    *string
		expectedError error
	}{
		{name: "no error while getting identity",
			mocked: mockSTS{
				gcio: &sts.GetCallerIdentityOutput{
					Account: aws.String("1"),
				},
				gcierr: nil,
			},
			expectedID:    aws.String("1"),
			expectedError: errors.New("at least one region name is required"),
		},
		{name: "error while getting identity",
			mocked: mockSTS{
				gcio:   nil,
				gcierr: errors.New("fail"),
			},
			expectedID:    nil,
			expectedError: errors.New("fail"),
		}}

	var ctx = context.Background()

	for i, tt := range tests {
		c := &connector{}
		err := c.setAccountID(ctx, tt.mocked)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(c.accountID, tt.expectedID) {
			t.Errorf("%s [%d] - accountID: received=%+v | expected=%+v",
				tt.name, i, c.accountID, tt.expectedID)
		}
	}
}

func TestConnector_setServices(t *testing.T) {
	tests := []struct {
		name         string
		configInput  *aws.Config
		credsInput   *credentials.Credentials
		regionsInput []string
		expectedSvcs []*serviceConnector
	}{
		{name: "one region service - no config",
			configInput:  nil,
			credsInput:   nil,
			regionsInput: []string{"eu-west-1"},
			expectedSvcs: []*serviceConnector{
				{region: "eu-west-1"},
			},
		},
		{name: "one region service - with config",
			configInput:  &aws.Config{},
			credsInput:   &credentials.Credentials{},
			regionsInput: []string{"eu-west-1"},
			expectedSvcs: []*serviceConnector{
				{region: "eu-west-1"},
			},
		},
		{name: "multiple region service - no config",
			configInput:  nil,
			regionsInput: []string{"eu-west-1", "eu-west-2", "eu-central-1"},
			expectedSvcs: []*serviceConnector{
				{region: "eu-west-1"},
				{region: "eu-west-2"},
				{region: "eu-central-1"},
			},
		},
		{name: "multiple region service - with config",
			configInput:  &aws.Config{},
			credsInput:   &credentials.Credentials{},
			regionsInput: []string{"eu-west-1", "eu-west-2", "eu-central-1"},
			expectedSvcs: []*serviceConnector{
				{region: "eu-west-1"},
				{region: "eu-west-2"},
				{region: "eu-central-1"},
			},
		}}

	for i, tt := range tests {
		c := &connector{
			regions: tt.regionsInput,
			creds:   tt.credsInput,
		}
		c.setServices(tt.configInput)
		for index, svc := range c.svcs {
			if svc.region != tt.expectedSvcs[index].region {
				t.Errorf("%s [%d] - services: received=%+v | expected=%+v",
					tt.name, i, c.svcs, tt.expectedSvcs)
			}
		}
	}
}
