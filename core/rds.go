package core

import (
	"github.com/aws/aws-sdk-go/service/rds"
)

// Returns all DB instances based on the input given
func (c *Connector) GetDBInstances(input *rds.DescribeDBInstancesInput) ([]*rds.DescribeDBInstancesOutput, error) {
	var errs = RawsErr{}
	var instances []*rds.DescribeDBInstancesOutput

	for _, svc := range c.svcs {
		if svc.rds == nil {
			svc.rds = rds.New(svc.session)
		}
		elb, err := svc.rds.DescribeDBInstances(input)
		instances = append(instances, elb)
		errs.AppendError(svc.region, rds.ServiceName, err)
	}
	return instances, errs
}

// Returns a list of tags from an ARN, extra filters for tags can also be provided
// For more information, please see: https://docs.aws.amazon.com/sdk-for-go/api/service/rds/#Filter
func (c *Connector) GetDBInstancesTags(input *rds.ListTagsForResourceInput) ([]*rds.ListTagsForResourceOutput, error) {
	var errs RawsErr = RawsErr{}
	var rdsTags []*rds.ListTagsForResourceOutput

	for _, svc := range c.svcs {
		if svc.rds == nil {
			svc.rds = rds.New(svc.session)
		}
		rdsTag, err := svc.rds.ListTagsForResource(input)
		rdsTags = append(rdsTags, rdsTag)
		errs.AppendError(svc.region, rds.ServiceName, err)
	}
	return rdsTags, errs
}
