package raws

import (
	"errors"
	"testing"

	"reflect"

	"github.com/aws/aws-sdk-go/aws"
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

func (m mockELB) DescribeLoadBalancers(input *elb.DescribeLoadBalancersInput) (*elb.DescribeLoadBalancersOutput, error) {
	return m.dlbo, m.dlberr
}

func (m mockELB) DescribeTags(input *elb.DescribeTagsInput) (*elb.DescribeTagsOutput, error) {
	return m.dto, m.dterr
}

func TestGetLoadBalancers(t *testing.T) {
	tests := []struct {
		name          string
		mocked        []*serviceConnector
		expectedELBs  []*elb.DescribeLoadBalancersOutput
		expectedError Errs
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
		expectedError: Errs{&callErr{
			err:     errors.New("error with test"),
			region:  "test",
			service: elb.ServiceName,
		}},
		expectedELBs: nil,
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
			expectedELBs: []*elb.DescribeLoadBalancersOutput{
				{
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
			expectedELBs: []*elb.DescribeLoadBalancersOutput{
				{
					LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
						{
							LoadBalancerName: aws.String("1"),
						},
					},
				},
				{
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
			expectedError: Errs{
				&callErr{
					err:     errors.New("error with test-1"),
					region:  "test-1",
					service: elb.ServiceName,
				},
			},
			expectedELBs: []*elb.DescribeLoadBalancersOutput{
				{
					LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
						{
							LoadBalancerName: aws.String("2"),
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		elbs, err := c.GetLoadBalancers(nil)
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
		expectedTags  []*elb.DescribeTagsOutput
		expectedError Errs
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
		expectedError: Errs{&callErr{
			err:     errors.New("error with test"),
			region:  "test",
			service: elb.ServiceName,
		}},
		expectedTags: nil,
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
			expectedTags: []*elb.DescribeTagsOutput{
				{
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
			expectedTags: []*elb.DescribeTagsOutput{
				{
					TagDescriptions: []*elb.TagDescription{
						{
							LoadBalancerName: aws.String("1"),
						},
					},
				},
				{
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
			expectedError: Errs{
				&callErr{
					err:     errors.New("error with test-1"),
					region:  "test-1",
					service: elb.ServiceName,
				},
			},
			expectedTags: []*elb.DescribeTagsOutput{
				{
					TagDescriptions: []*elb.TagDescription{
						{
							LoadBalancerName: aws.String("2"),
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		tags, err := c.GetLoadBalancersTags(nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(tags, tt.expectedTags) {
			t.Errorf("%s [%d] - tags: received=%+v | expected=%+v",
				tt.name, i, tags, tt.expectedTags)
		}
	}
}
