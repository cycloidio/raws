package raws

import (
	"context"
	"errors"
	"testing"

	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
)

type mockELB struct {
	elbiface.ELBAPI

	// Mocking of DescribeLoadBalancers
	dlbo   *elb.DescribeLoadBalancersOutput
	dlberr error

	// Mocking of DescribeTags
	dto   *elb.DescribeTagsOutput
	dterr error
}

func (m mockELB) DescribeLoadBalancersWithContext(
	_ aws.Context, _ *elb.DescribeLoadBalancersInput, _ ...request.Option,
) (*elb.DescribeLoadBalancersOutput, error) {
	return m.dlbo, m.dlberr
}

func (m mockELB) DescribeTagsWithContext(
	_ aws.Context, _ *elb.DescribeTagsInput, _ ...request.Option,
) (*elb.DescribeTagsOutput, error) {
	return m.dto, m.dterr
}

func TestGetLoadBalancers(t *testing.T) {
	tests := []struct {
		name          string
		mocked        []*serviceConnector
		expectedELBs  map[string]elb.DescribeLoadBalancersOutput
		expectedError error
	}{{name: "one region no error",
		mocked: []*serviceConnector{
			{
				region: "test",
				elb: mockELB{
					dlbo:   &elb.DescribeLoadBalancersOutput{},
					dlberr: errors.New("error with test"),
				},
			},
		},
		expectedError: Errors{Error{
			err:     errors.New("error with test"),
			region:  "test",
			service: elb.ServiceName,
		}},
		expectedELBs: map[string]elb.DescribeLoadBalancersOutput{},
	},
		{name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					elb: mockELB{
						dlbo: &elb.DescribeLoadBalancersOutput{
							LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
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
			expectedELBs: map[string]elb.DescribeLoadBalancersOutput{
				"test": {
					LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
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
					elb: mockELB{
						dlbo: &elb.DescribeLoadBalancersOutput{
							LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
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
					elb: mockELB{
						dlbo: &elb.DescribeLoadBalancersOutput{
							LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
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
			expectedELBs: map[string]elb.DescribeLoadBalancersOutput{
				"test-1": {
					LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
						{
							LoadBalancerName: aws.String("1"),
						},
					},
				},
				"test-2": {
					LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
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
					elb: mockELB{
						dlbo:   &elb.DescribeLoadBalancersOutput{},
						dlberr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					elb: mockELB{
						dlbo: &elb.DescribeLoadBalancersOutput{
							LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
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
					service: elb.ServiceName,
				},
			},
			expectedELBs: map[string]elb.DescribeLoadBalancersOutput{
				"test-2": {
					LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
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
		elbs, err := c.GetLoadBalancers(ctx, nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(elbs, tt.expectedELBs) {
			t.Errorf("%s [%d] - ELBs (v1): received=%+v | expected=%+v",
				tt.name, i, elbs, tt.expectedELBs)
		}
	}
}

func TestGetLoadBalancersTags(t *testing.T) {
	tests := []struct {
		name          string
		mocked        []*serviceConnector
		expectedTags  map[string]elb.DescribeTagsOutput
		expectedError error
	}{{name: "one region with error",
		mocked: []*serviceConnector{
			{
				region: "test",
				elb: mockELB{
					dto:   &elb.DescribeTagsOutput{},
					dterr: errors.New("error with test"),
				},
			},
		},
		expectedError: Errors{Error{
			err:     errors.New("error with test"),
			region:  "test",
			service: elb.ServiceName,
		}},
		expectedTags: map[string]elb.DescribeTagsOutput{},
	},
		{name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					elb: mockELB{
						dto: &elb.DescribeTagsOutput{
							TagDescriptions: []*elb.TagDescription{
								{
									LoadBalancerName: aws.String("1"),
								},
							},
						},
						dterr: nil,
					},
				},
			},
			expectedError: nil,
			expectedTags: map[string]elb.DescribeTagsOutput{
				"test": {
					TagDescriptions: []*elb.TagDescription{
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
					elb: mockELB{
						dto: &elb.DescribeTagsOutput{
							TagDescriptions: []*elb.TagDescription{
								{
									LoadBalancerName: aws.String("1"),
								},
							},
						},
						dterr: nil,
					},
				},
				{
					region: "test-2",
					elb: mockELB{
						dto: &elb.DescribeTagsOutput{
							TagDescriptions: []*elb.TagDescription{
								{
									LoadBalancerName: aws.String("2"),
								},
							},
						},
						dterr: nil,
					},
				},
			},
			expectedError: nil,
			expectedTags: map[string]elb.DescribeTagsOutput{
				"test-1": {
					TagDescriptions: []*elb.TagDescription{
						{
							LoadBalancerName: aws.String("1"),
						},
					},
				},
				"test-2": {
					TagDescriptions: []*elb.TagDescription{
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
					elb: mockELB{
						dto:   &elb.DescribeTagsOutput{},
						dterr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					elb: mockELB{
						dto: &elb.DescribeTagsOutput{
							TagDescriptions: []*elb.TagDescription{
								{
									LoadBalancerName: aws.String("2"),
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
					service: elb.ServiceName,
				},
			},
			expectedTags: map[string]elb.DescribeTagsOutput{
				"test-2": {
					TagDescriptions: []*elb.TagDescription{
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
		tags, err := c.GetLoadBalancersTags(ctx, nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(tags, tt.expectedTags) {
			t.Errorf("%s [%d] - tags: received=%+v | expected=%+v",
				tt.name, i, tags, tt.expectedTags)
		}
	}
}
