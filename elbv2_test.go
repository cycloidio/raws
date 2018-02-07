package raws

import (
	"context"
	"errors"
	"testing"

	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
)

type mockELBV2 struct {
	elbv2iface.ELBV2API

	// Mocking of DescribeLoadBalancers
	dlbo   *elbv2.DescribeLoadBalancersOutput
	dlberr error

	// Mocking of DescribeTags
	dto   *elbv2.DescribeTagsOutput
	dterr error
}

func (m mockELBV2) DescribeLoadBalancersWithContext(
	_ aws.Context, _ *elbv2.DescribeLoadBalancersInput, _ ...request.Option,
) (*elbv2.DescribeLoadBalancersOutput, error) {
	return m.dlbo, m.dlberr
}

func (m mockELBV2) DescribeTagsWithContext(
	_ aws.Context, _ *elbv2.DescribeTagsInput, _ ...request.Option,
) (*elbv2.DescribeTagsOutput, error) {
	return m.dto, m.dterr
}

func TestGetLoadBalancersV2(t *testing.T) {
	tests := []struct {
		name          string
		mocked        []*serviceConnector
		expectedELBs  map[string]elbv2.DescribeLoadBalancersOutput
		expectedError error
	}{{name: "one region with error",
		mocked: []*serviceConnector{
			{
				region: "test",
				elbv2: mockELBV2{
					dlbo:   &elbv2.DescribeLoadBalancersOutput{},
					dlberr: errors.New("error with test"),
				},
			},
		},
		expectedError: Errors{Error{
			err:     errors.New("error with test"),
			region:  "test",
			service: elbv2.ServiceName,
		}},
		expectedELBs: map[string]elbv2.DescribeLoadBalancersOutput{},
	},
		{name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					elbv2: mockELBV2{
						dlbo: &elbv2.DescribeLoadBalancersOutput{
							LoadBalancers: []*elbv2.LoadBalancer{
								{
									LoadBalancerName: aws.String("1"),
								},
							},
						},
						dlberr: nil,
					},
				},
			},
			expectedError: nil,
			expectedELBs: map[string]elbv2.DescribeLoadBalancersOutput{
				"test": {
					LoadBalancers: []*elbv2.LoadBalancer{
						{
							LoadBalancerName: aws.String("1"),
						},
					},
				},
			},
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					elbv2: mockELBV2{
						dlbo: &elbv2.DescribeLoadBalancersOutput{
							LoadBalancers: []*elbv2.LoadBalancer{
								{
									LoadBalancerName: aws.String("1"),
								},
							},
						},
						dlberr: nil,
					},
				},
				{
					region: "test-2",
					elbv2: mockELBV2{
						dlbo: &elbv2.DescribeLoadBalancersOutput{
							LoadBalancers: []*elbv2.LoadBalancer{
								{
									LoadBalancerName: aws.String("2"),
								},
							},
						},
						dlberr: nil,
					},
				},
			},
			expectedError: nil,
			expectedELBs: map[string]elbv2.DescribeLoadBalancersOutput{
				"test-1": {
					LoadBalancers: []*elbv2.LoadBalancer{
						{
							LoadBalancerName: aws.String("1"),
						},
					},
				},
				"test-2": {
					LoadBalancers: []*elbv2.LoadBalancer{
						{
							LoadBalancerName: aws.String("2"),
						},
					},
				},
			},
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					elbv2: mockELBV2{
						dlbo:   &elbv2.DescribeLoadBalancersOutput{},
						dlberr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					elbv2: mockELBV2{
						dlbo: &elbv2.DescribeLoadBalancersOutput{
							LoadBalancers: []*elbv2.LoadBalancer{
								{
									LoadBalancerName: aws.String("2"),
								},
							},
						},
						dlberr: nil,
					},
				},
			},
			expectedError: Errors{
				Error{
					err:     errors.New("error with test-1"),
					region:  "test-1",
					service: elbv2.ServiceName,
				},
			},
			expectedELBs: map[string]elbv2.DescribeLoadBalancersOutput{
				"test-2": {
					LoadBalancers: []*elbv2.LoadBalancer{
						{
							LoadBalancerName: aws.String("2"),
						},
					},
				},
			},
		},
	}

	var ctx = context.Background()

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		elbs, err := c.GetLoadBalancersV2(ctx, nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(elbs, tt.expectedELBs) {
			t.Errorf("%s [%d] - ELBs (v1): received=%+v | expected=%+v",
				tt.name, i, elbs, tt.expectedELBs)
		}
	}
}

func TestGetLoadBalancersV2Tags(t *testing.T) {
	tests := []struct {
		name          string
		mocked        []*serviceConnector
		expectedTags  map[string]elbv2.DescribeTagsOutput
		expectedError error
	}{{name: "one region with error",
		mocked: []*serviceConnector{
			{
				region: "test",
				elbv2: mockELBV2{
					dto:   &elbv2.DescribeTagsOutput{},
					dterr: errors.New("error with test"),
				},
			},
		},
		expectedError: Errors{Error{
			err:     errors.New("error with test"),
			region:  "test",
			service: elbv2.ServiceName,
		}},
		expectedTags: map[string]elbv2.DescribeTagsOutput{},
	},
		{name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					elbv2: mockELBV2{
						dto: &elbv2.DescribeTagsOutput{
							TagDescriptions: []*elbv2.TagDescription{
								{
									ResourceArn: aws.String("1"),
								},
							},
						},
						dterr: nil,
					},
				},
			},
			expectedError: nil,
			expectedTags: map[string]elbv2.DescribeTagsOutput{
				"test": {
					TagDescriptions: []*elbv2.TagDescription{
						{
							ResourceArn: aws.String("1"),
						},
					},
				},
			},
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					elbv2: mockELBV2{
						dto: &elbv2.DescribeTagsOutput{
							TagDescriptions: []*elbv2.TagDescription{
								{
									ResourceArn: aws.String("1"),
								},
							},
						},
						dterr: nil,
					},
				},
				{
					region: "test-2",
					elbv2: mockELBV2{
						dto: &elbv2.DescribeTagsOutput{
							TagDescriptions: []*elbv2.TagDescription{
								{
									ResourceArn: aws.String("2"),
								},
							},
						},
						dterr: nil,
					},
				},
			},
			expectedError: nil,
			expectedTags: map[string]elbv2.DescribeTagsOutput{
				"test-1": {
					TagDescriptions: []*elbv2.TagDescription{
						{
							ResourceArn: aws.String("1"),
						},
					},
				},
				"test-2": {
					TagDescriptions: []*elbv2.TagDescription{
						{
							ResourceArn: aws.String("2"),
						},
					},
				},
			},
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					elbv2: mockELBV2{
						dto:   &elbv2.DescribeTagsOutput{},
						dterr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					elbv2: mockELBV2{
						dto: &elbv2.DescribeTagsOutput{
							TagDescriptions: []*elbv2.TagDescription{
								{
									ResourceArn: aws.String("2"),
								},
							},
						},
						dterr: nil,
					},
				},
			},
			expectedError: Errors{
				Error{
					err:     errors.New("error with test-1"),
					region:  "test-1",
					service: elbv2.ServiceName,
				},
			},
			expectedTags: map[string]elbv2.DescribeTagsOutput{
				"test-2": {
					TagDescriptions: []*elbv2.TagDescription{
						{
							ResourceArn: aws.String("2"),
						},
					},
				},
			},
		},
	}

	var ctx = context.Background()

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		tags, err := c.GetLoadBalancersV2Tags(ctx, nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(tags, tt.expectedTags) {
			t.Errorf("%s [%d] - tags: received=%+v | expected=%+v",
				tt.name, i, tags, tt.expectedTags)
		}
	}
}
