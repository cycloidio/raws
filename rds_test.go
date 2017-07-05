package raws

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

type mockRDS struct {
	rdsiface.RDSAPI

	// Mocking of DescribeDBInstances
	ddio   *rds.DescribeDBInstancesOutput
	ddierr error

	// Mocking of ListTagsForResource
	ltfro   *rds.ListTagsForResourceOutput
	ltfrerr error
}

func (m mockRDS) DescribeDBInstances(input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error) {
	return m.ddio, m.ddierr
}

func (m mockRDS) ListTagsForResource(input *rds.ListTagsForResourceInput) (*rds.ListTagsForResourceOutput, error) {
	return m.ltfro, m.ltfrerr
}

func TestGetDBInstances(t *testing.T) {
	tests := []struct {
		name              string
		mocked            []*serviceConnector
		expectedInstances []*rds.DescribeDBInstancesOutput
		expectedError     Errs
	}{{
		name: "one region no error",
		mocked: []*serviceConnector{
			{
				region: "test",
				rds: mockRDS{
					ddio: &rds.DescribeDBInstancesOutput{
						DBInstances: []*rds.DBInstance{
							{
								DbiResourceId: aws.String("1"),
							},
						},
					},
					ddierr: nil,
				},
			},
		},
		expectedError: nil,
		expectedInstances: []*rds.DescribeDBInstancesOutput{
			{
				DBInstances: []*rds.DBInstance{
					{
						DbiResourceId: aws.String("1"),
					},
				},
			},
		},
	},
		{name: "one region with error",
			mocked: []*serviceConnector{
				{
					region: "test",
					rds: mockRDS{
						ddio:   nil,
						ddierr: errors.New("error with test"),
					},
				},
			},
			expectedError: Errs{
				&callErr{
					err:     errors.New("error with test"),
					region:  "test",
					service: rds.ServiceName,
				},
			},
			expectedInstances: nil,
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					rds: mockRDS{
						ddio:   &rds.DescribeDBInstancesOutput{},
						ddierr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					rds: mockRDS{
						ddio: &rds.DescribeDBInstancesOutput{
							DBInstances: []*rds.DBInstance{
								{
									DbiResourceId: aws.String("2"),
								},
							},
						},
						ddierr: nil,
					},
				},
			},
			expectedError: Errs{
				&callErr{
					err:     errors.New("error with test-1"),
					region:  "test-1",
					service: rds.ServiceName,
				},
			},
			expectedInstances: []*rds.DescribeDBInstancesOutput{
				{
					DBInstances: []*rds.DBInstance{
						{
							DbiResourceId: aws.String("2"),
						},
					},
				},
			},
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					rds: mockRDS{
						ddio: &rds.DescribeDBInstancesOutput{
							DBInstances: []*rds.DBInstance{
								{
									DbiResourceId: aws.String("1"),
								},
							},
						},
						ddierr: nil,
					},
				},
				{
					region: "test-2",
					rds: mockRDS{
						ddio: &rds.DescribeDBInstancesOutput{
							DBInstances: []*rds.DBInstance{
								{
									DbiResourceId: aws.String("2"),
								},
							},
						},
						ddierr: nil,
					},
				},
			},
			expectedError: nil,
			expectedInstances: []*rds.DescribeDBInstancesOutput{
				{
					DBInstances: []*rds.DBInstance{
						{
							DbiResourceId: aws.String("1"),
						},
					},
				},
				{
					DBInstances: []*rds.DBInstance{
						{
							DbiResourceId: aws.String("2"),
						},
					},
				},
			},
		}}

	for i, tt := range tests {
		c := &Connector{svcs: tt.mocked}
		instances, err := c.GetDBInstances(nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(instances, tt.expectedInstances) {
			t.Errorf("%s [%d] - DB instances: received=%+v | expected=%+v",
				tt.name, i, instances, tt.expectedInstances)
		}
	}
}

func TestGetDBInstancesTags(t *testing.T) {
	tests := []struct {
		name          string
		mocked        []*serviceConnector
		expectedTags  []*rds.ListTagsForResourceOutput
		expectedError Errs
	}{
		{name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					rds: mockRDS{
						ltfro: &rds.ListTagsForResourceOutput{
							TagList: []*rds.Tag{
								{
									Key:   aws.String("test"),
									Value: aws.String("1"),
								},
							},
						},
						ltfrerr: nil,
					},
				},
			},
			expectedError: nil,
			expectedTags: []*rds.ListTagsForResourceOutput{
				{
					TagList: []*rds.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("1"),
						},
					},
				},
			},
		},
		{name: "one region with error",
			mocked: []*serviceConnector{
				{
					region: "test",
					rds: mockRDS{
						ltfro:   nil,
						ltfrerr: errors.New("error with test"),
					},
				},
			},
			expectedError: Errs{
				&callErr{
					err:     errors.New("error with test"),
					region:  "test",
					service: rds.ServiceName,
				},
			},
			expectedTags: nil,
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					rds: mockRDS{
						ltfro: &rds.ListTagsForResourceOutput{
							TagList: []*rds.Tag{
								{
									Key:   aws.String("test"),
									Value: aws.String("1"),
								},
							},
						},
						ltfrerr: nil,
					},
				},
				{
					region: "test-2",
					rds: mockRDS{
						ltfro: &rds.ListTagsForResourceOutput{
							TagList: []*rds.Tag{
								{
									Key:   aws.String("test"),
									Value: aws.String("2"),
								},
							},
						},
						ltfrerr: nil,
					},
				},
			},
			expectedError: nil,
			expectedTags: []*rds.ListTagsForResourceOutput{
				{
					TagList: []*rds.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("1"),
						},
					},
				},
				{
					TagList: []*rds.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("2"),
						},
					},
				},
			},
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					rds: mockRDS{
						ltfro:   nil,
						ltfrerr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					rds: mockRDS{
						ltfro: &rds.ListTagsForResourceOutput{
							TagList: []*rds.Tag{
								{
									Key:   aws.String("test"),
									Value: aws.String("2"),
								},
							},
						},
						ltfrerr: nil,
					},
				},
			},
			expectedError: Errs{
				&callErr{
					err:     errors.New("error with test-1"),
					region:  "test-1",
					service: rds.ServiceName,
				},
			},
			expectedTags: []*rds.ListTagsForResourceOutput{
				{
					TagList: []*rds.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("2"),
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		c := &Connector{svcs: tt.mocked}
		tags, err := c.GetDBInstancesTags(nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(tags, tt.expectedTags) {
			t.Errorf("%s [%d] - DB instances: received=%+v | expected=%+v",
				tt.name, i, tags, tt.expectedTags)
		}
	}
}
