package core

import (
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

func (m mockRDS) DescribeIDBInstances(input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error) {
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
}
