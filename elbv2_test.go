package raws

import (
	"errors"
	"testing"

	"reflect"

	"github.com/aws/aws-sdk-go/aws"
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

func (m mockELBV2) DescribeLoadBalancers(input *elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
	return m.dlbo, m.dlberr
}

func (m mockELBV2) DescribeTags(input *elbv2.DescribeTagsInput) (*elbv2.DescribeTagsOutput, error) {
	return m.dto, m.dterr
}

func TestGetLoadBalancersV2(t *testing.T) {
	tests := []struct {
		name          string
		mocked        []*serviceConnector
		expectedELBs  []*elbv2.DescribeLoadBalancersOutput
		expectedError Errs
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
		expectedError: Errs{&callErr{
			err:     errors.New("error with test"),
			region:  "test",
			service: elbv2.ServiceName,
		}},
		expectedELBs: nil,
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
			expectedELBs: []*elbv2.DescribeLoadBalancersOutput{
				{
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
			expectedELBs: []*elbv2.DescribeLoadBalancersOutput{
				{
					LoadBalancers: []*elbv2.LoadBalancer{
						{
							LoadBalancerName: aws.String("1"),
						},
					},
				},
				{
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
			expectedError: Errs{
				&callErr{
					err:     errors.New("error with test-1"),
					region:  "test-1",
					service: elbv2.ServiceName,
				},
			},
			expectedELBs: []*elbv2.DescribeLoadBalancersOutput{
				{
					LoadBalancers: []*elbv2.LoadBalancer{
						{
							LoadBalancerName: aws.String("2"),
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		c := &Connector{svcs: tt.mocked}
		elbs, err := c.GetLoadBalancersV2(nil)
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
		expectedTags  []*elbv2.DescribeTagsOutput
		expectedError Errs
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
		expectedError: Errs{&callErr{
			err:     errors.New("error with test"),
			region:  "test",
			service: elbv2.ServiceName,
		}},
		expectedTags: nil,
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
			expectedTags: []*elbv2.DescribeTagsOutput{
				{
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
			expectedTags: []*elbv2.DescribeTagsOutput{
				{
					TagDescriptions: []*elbv2.TagDescription{
						{
							ResourceArn: aws.String("1"),
						},
					},
				},
				{
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
			expectedError: Errs{
				&callErr{
					err:     errors.New("error with test-1"),
					region:  "test-1",
					service: elbv2.ServiceName,
				},
			},
			expectedTags: []*elbv2.DescribeTagsOutput{
				{
					TagDescriptions: []*elbv2.TagDescription{
						{
							ResourceArn: aws.String("2"),
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		c := &Connector{svcs: tt.mocked}
		tags, err := c.GetLoadBalancersV2Tags(nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(tags, tt.expectedTags) {
			t.Errorf("%s [%d] - tags: received=%+v | expected=%+v",
				tt.name, i, tags, tt.expectedTags)
		}
	}
}
