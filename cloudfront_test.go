package raws

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
)

type mockCloudFront struct {
	cloudfrontiface.CloudFrontAPI

	// Mock of DescribeInstances
	ldo    *cloudfront.ListDistributionsOutput
	ldoerr error
}

func (m mockCloudFront) ListDistributionsWithContext(
	_ aws.Context, _ *cloudfront.ListDistributionsInput, _ ...request.Option,
) (*cloudfront.ListDistributionsOutput, error) {
	return m.ldo, m.ldoerr
}

func TestGetCloudFrontDistributions(t *testing.T) {
	tests := []struct {
		name        string
		mocked      []*serviceConnector
		expectedOpt map[string]cloudfront.ListDistributionsOutput
		expectedErr error
	}{
		{
			name: "one region with error",
			mocked: []*serviceConnector{
				{
					region: "test",
					cloudfront: mockCloudFront{
						ldo:    &cloudfront.ListDistributionsOutput{},
						ldoerr: errors.New("error with test"),
					},
				},
			},
			expectedErr: Errors{Error{
				err:     errors.New("error with test"),
				region:  "test",
				service: cloudfront.ServiceName,
			}},
			expectedOpt: map[string]cloudfront.ListDistributionsOutput{},
		},
		{
			name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					cloudfront: mockCloudFront{
						ldo: &cloudfront.ListDistributionsOutput{
							DistributionList: &cloudfront.DistributionList{
								Items: []*cloudfront.DistributionSummary{
									&cloudfront.DistributionSummary{
										Id: aws.String("123"),
									},
								},
							},
						},
						ldoerr: nil,
					},
				},
			},
			expectedErr: nil,
			expectedOpt: map[string]cloudfront.ListDistributionsOutput{
				"test": {
					DistributionList: &cloudfront.DistributionList{
						Items: []*cloudfront.DistributionSummary{
							&cloudfront.DistributionSummary{
								Id: aws.String("123"),
							},
						},
					},
				},
			},
		},
		{
			name: "multiple region no error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					cloudfront: mockCloudFront{
						ldo: &cloudfront.ListDistributionsOutput{
							DistributionList: &cloudfront.DistributionList{
								Items: []*cloudfront.DistributionSummary{
									&cloudfront.DistributionSummary{
										Id: aws.String("123"),
									},
								},
							},
						},
						ldoerr: nil,
					},
				},
				{
					region: "test-2",
					cloudfront: mockCloudFront{
						ldo: &cloudfront.ListDistributionsOutput{
							DistributionList: &cloudfront.DistributionList{
								Items: []*cloudfront.DistributionSummary{
									&cloudfront.DistributionSummary{
										Id: aws.String("456"),
									},
								},
							},
						},
						ldoerr: nil,
					},
				},
			},
			expectedErr: nil,
			expectedOpt: map[string]cloudfront.ListDistributionsOutput{
				"test-1": {
					DistributionList: &cloudfront.DistributionList{
						Items: []*cloudfront.DistributionSummary{
							&cloudfront.DistributionSummary{
								Id: aws.String("123"),
							},
						},
					},
				},
				"test-2": {
					DistributionList: &cloudfront.DistributionList{
						Items: []*cloudfront.DistributionSummary{
							&cloudfront.DistributionSummary{
								Id: aws.String("456"),
							},
						},
					},
				},
			},
		},
		{
			name: "multiple region with error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					cloudfront: mockCloudFront{
						ldo:    &cloudfront.ListDistributionsOutput{},
						ldoerr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					cloudfront: mockCloudFront{
						ldo: &cloudfront.ListDistributionsOutput{
							DistributionList: &cloudfront.DistributionList{
								Items: []*cloudfront.DistributionSummary{
									&cloudfront.DistributionSummary{
										Id: aws.String("456"),
									},
								},
							},
						},
						ldoerr: nil,
					},
				},
			},
			expectedErr: Errors{Error{
				err:     errors.New("error with test-1"),
				region:  "test-1",
				service: cloudfront.ServiceName,
			}},
			expectedOpt: map[string]cloudfront.ListDistributionsOutput{
				"test-2": {
					DistributionList: &cloudfront.DistributionList{
						Items: []*cloudfront.DistributionSummary{
							&cloudfront.DistributionSummary{
								Id: aws.String("456"),
							},
						},
					},
				},
			},
		},
	}

	var ctx = context.Background()

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		opt, err := c.GetCloudFrontDistributions(ctx, nil)
		checkErrors(t, tt.name, i, err, tt.expectedErr)
		if !reflect.DeepEqual(opt, tt.expectedOpt) {
			t.Errorf("%s [%d] - EC2 instances: received=%+v | expected=%+v",
				tt.name, i, opt, tt.expectedOpt)
		}
	}
}
