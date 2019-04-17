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

	// Mock of ListDistributionsOutput
	ldo   *cloudfront.ListDistributionsOutput
	lderr error

	// Mock of  ListPublicKeysOutput
	lpko   *cloudfront.ListPublicKeysOutput
	lpkerr error

	loaio   *cloudfront.ListCloudFrontOriginAccessIdentitiesOutput
	loaierr error
}

func (m mockCloudFront) ListDistributionsWithContext(
	_ aws.Context, _ *cloudfront.ListDistributionsInput, _ ...request.Option,
) (*cloudfront.ListDistributionsOutput, error) {
	return m.ldo, m.lderr
}

func (m mockCloudFront) ListPublicKeysWithContext(
	_ aws.Context, _ *cloudfront.ListPublicKeysInput, _ ...request.Option,
) (*cloudfront.ListPublicKeysOutput, error) {
	return m.lpko, m.lpkerr
}

func (m mockCloudFront) ListCloudFrontOriginAccessIdentitiesWithContext(
	_ aws.Context, _ *cloudfront.ListCloudFrontOriginAccessIdentitiesInput, _ ...request.Option,
) (*cloudfront.ListCloudFrontOriginAccessIdentitiesOutput, error) {
	return m.loaio, m.loaierr
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
						ldo:   &cloudfront.ListDistributionsOutput{},
						lderr: errors.New("error with test"),
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
						lderr: nil,
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
						lderr: nil,
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
						lderr: nil,
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
						ldo:   &cloudfront.ListDistributionsOutput{},
						lderr: errors.New("error with test-1"),
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
						lderr: nil,
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

func TestGetCloudFrontPublicKeys(t *testing.T) {
	tests := []struct {
		name        string
		mocked      []*serviceConnector
		expectedOpt map[string]cloudfront.ListPublicKeysOutput
		expectedErr error
	}{
		{
			name: "one region with error",
			mocked: []*serviceConnector{
				{
					region: "test",
					cloudfront: mockCloudFront{
						lpko:   &cloudfront.ListPublicKeysOutput{},
						lpkerr: errors.New("error with test"),
					},
				},
			},
			expectedErr: Errors{Error{
				err:     errors.New("error with test"),
				region:  "test",
				service: cloudfront.ServiceName,
			}},
			expectedOpt: map[string]cloudfront.ListPublicKeysOutput{},
		},
		{
			name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					cloudfront: mockCloudFront{
						lpko: &cloudfront.ListPublicKeysOutput{
							PublicKeyList: &cloudfront.PublicKeyList{
								Items: []*cloudfront.PublicKeySummary{
									&cloudfront.PublicKeySummary{
										Id: aws.String("123"),
									},
								},
							},
						},
						lpkerr: nil,
					},
				},
			},
			expectedErr: nil,
			expectedOpt: map[string]cloudfront.ListPublicKeysOutput{
				"test": {
					PublicKeyList: &cloudfront.PublicKeyList{
						Items: []*cloudfront.PublicKeySummary{
							&cloudfront.PublicKeySummary{
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
						lpko: &cloudfront.ListPublicKeysOutput{
							PublicKeyList: &cloudfront.PublicKeyList{
								Items: []*cloudfront.PublicKeySummary{
									&cloudfront.PublicKeySummary{
										Id: aws.String("123"),
									},
								},
							},
						},
						lpkerr: nil,
					},
				},
				{
					region: "test-2",
					cloudfront: mockCloudFront{
						lpko: &cloudfront.ListPublicKeysOutput{
							PublicKeyList: &cloudfront.PublicKeyList{
								Items: []*cloudfront.PublicKeySummary{
									&cloudfront.PublicKeySummary{
										Id: aws.String("456"),
									},
								},
							},
						},
						lpkerr: nil,
					},
				},
			},
			expectedErr: nil,
			expectedOpt: map[string]cloudfront.ListPublicKeysOutput{
				"test-1": {
					PublicKeyList: &cloudfront.PublicKeyList{
						Items: []*cloudfront.PublicKeySummary{
							&cloudfront.PublicKeySummary{
								Id: aws.String("123"),
							},
						},
					},
				},
				"test-2": {
					PublicKeyList: &cloudfront.PublicKeyList{
						Items: []*cloudfront.PublicKeySummary{
							&cloudfront.PublicKeySummary{
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
						lpko:   &cloudfront.ListPublicKeysOutput{},
						lpkerr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					cloudfront: mockCloudFront{
						lpko: &cloudfront.ListPublicKeysOutput{
							PublicKeyList: &cloudfront.PublicKeyList{
								Items: []*cloudfront.PublicKeySummary{
									&cloudfront.PublicKeySummary{
										Id: aws.String("456"),
									},
								},
							},
						},
						lpkerr: nil,
					},
				},
			},
			expectedErr: Errors{Error{
				err:     errors.New("error with test-1"),
				region:  "test-1",
				service: cloudfront.ServiceName,
			}},
			expectedOpt: map[string]cloudfront.ListPublicKeysOutput{
				"test-2": {
					PublicKeyList: &cloudfront.PublicKeyList{
						Items: []*cloudfront.PublicKeySummary{
							&cloudfront.PublicKeySummary{
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
		opt, err := c.GetCloudFrontPublicKeys(ctx, nil)
		checkErrors(t, tt.name, i, err, tt.expectedErr)
		if !reflect.DeepEqual(opt, tt.expectedOpt) {
			t.Errorf("%s [%d] - EC2 instances: received=%+v | expected=%+v",
				tt.name, i, opt, tt.expectedOpt)
		}
	}
}

func TestGetCloudFrontOriginAccessIdentities(t *testing.T) {
	tests := []struct {
		name        string
		mocked      []*serviceConnector
		expectedOpt map[string]cloudfront.ListCloudFrontOriginAccessIdentitiesOutput
		expectedErr error
	}{
		{
			name: "one region with error",
			mocked: []*serviceConnector{
				{
					region: "test",
					cloudfront: mockCloudFront{
						loaio:   &cloudfront.ListCloudFrontOriginAccessIdentitiesOutput{},
						loaierr: errors.New("error with test"),
					},
				},
			},
			expectedErr: Errors{Error{
				err:     errors.New("error with test"),
				region:  "test",
				service: cloudfront.ServiceName,
			}},
			expectedOpt: map[string]cloudfront.ListCloudFrontOriginAccessIdentitiesOutput{},
		},
		{
			name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					cloudfront: mockCloudFront{
						loaio: &cloudfront.ListCloudFrontOriginAccessIdentitiesOutput{
							CloudFrontOriginAccessIdentityList: &cloudfront.OriginAccessIdentityList{
								Items: []*cloudfront.OriginAccessIdentitySummary{
									&cloudfront.OriginAccessIdentitySummary{
										Id: aws.String("123"),
									},
								},
							},
						},
						loaierr: nil,
					},
				},
			},
			expectedErr: nil,
			expectedOpt: map[string]cloudfront.ListCloudFrontOriginAccessIdentitiesOutput{
				"test": {
					CloudFrontOriginAccessIdentityList: &cloudfront.OriginAccessIdentityList{
						Items: []*cloudfront.OriginAccessIdentitySummary{
							&cloudfront.OriginAccessIdentitySummary{
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
						loaio: &cloudfront.ListCloudFrontOriginAccessIdentitiesOutput{
							CloudFrontOriginAccessIdentityList: &cloudfront.OriginAccessIdentityList{
								Items: []*cloudfront.OriginAccessIdentitySummary{
									&cloudfront.OriginAccessIdentitySummary{
										Id: aws.String("123"),
									},
								},
							},
						},
						loaierr: nil,
					},
				},
				{
					region: "test-2",
					cloudfront: mockCloudFront{
						loaio: &cloudfront.ListCloudFrontOriginAccessIdentitiesOutput{
							CloudFrontOriginAccessIdentityList: &cloudfront.OriginAccessIdentityList{
								Items: []*cloudfront.OriginAccessIdentitySummary{
									&cloudfront.OriginAccessIdentitySummary{
										Id: aws.String("456"),
									},
								},
							},
						},
						loaierr: nil,
					},
				},
			},
			expectedErr: nil,
			expectedOpt: map[string]cloudfront.ListCloudFrontOriginAccessIdentitiesOutput{
				"test-1": {
					CloudFrontOriginAccessIdentityList: &cloudfront.OriginAccessIdentityList{
						Items: []*cloudfront.OriginAccessIdentitySummary{
							&cloudfront.OriginAccessIdentitySummary{
								Id: aws.String("123"),
							},
						},
					},
				},
				"test-2": {
					CloudFrontOriginAccessIdentityList: &cloudfront.OriginAccessIdentityList{
						Items: []*cloudfront.OriginAccessIdentitySummary{
							&cloudfront.OriginAccessIdentitySummary{
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
						loaio:   &cloudfront.ListCloudFrontOriginAccessIdentitiesOutput{},
						loaierr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					cloudfront: mockCloudFront{
						loaio: &cloudfront.ListCloudFrontOriginAccessIdentitiesOutput{
							CloudFrontOriginAccessIdentityList: &cloudfront.OriginAccessIdentityList{
								Items: []*cloudfront.OriginAccessIdentitySummary{
									&cloudfront.OriginAccessIdentitySummary{
										Id: aws.String("456"),
									},
								},
							},
						},
						loaierr: nil,
					},
				},
			},
			expectedErr: Errors{Error{
				err:     errors.New("error with test-1"),
				region:  "test-1",
				service: cloudfront.ServiceName,
			}},
			expectedOpt: map[string]cloudfront.ListCloudFrontOriginAccessIdentitiesOutput{
				"test-2": {
					CloudFrontOriginAccessIdentityList: &cloudfront.OriginAccessIdentityList{
						Items: []*cloudfront.OriginAccessIdentitySummary{
							&cloudfront.OriginAccessIdentitySummary{
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
		opt, err := c.GetCloudFrontOriginAccessIdentities(ctx, nil)
		checkErrors(t, tt.name, i, err, tt.expectedErr)
		if !reflect.DeepEqual(opt, tt.expectedOpt) {
			t.Errorf("%s [%d] - EC2 instances: received=%+v | expected=%+v",
				tt.name, i, opt, tt.expectedOpt)
		}
	}
}
