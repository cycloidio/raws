package core

import (
	"github.com/aws/aws-sdk-go/service/rds"
)

// Returns all DB instances based on the input given
func (c *Connector) GetDBInstances(input *rds.DescribeDBInstancesInput) ([]*rds.DescribeDBInstancesOutput, Errs) {
	var errs Errs
	var instances []*rds.DescribeDBInstancesOutput

	for _, svc := range c.svcs {
		if svc.rds == nil {
			svc.rds = rds.New(svc.session)
		}
		instance, err := svc.rds.DescribeDBInstances(input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, rds.ServiceName, err))
		} else {
			instances = append(instances, instance)
		}
	}
	return instances, errs
}

// Returns a list of tags from an ARN, extra filters for tags can also be provided
// For more information, please see: https://docs.aws.amazon.com/sdk-for-go/api/service/rds/#Filter
func (c *Connector) GetDBInstancesTags(input *rds.ListTagsForResourceInput) ([]*rds.ListTagsForResourceOutput, Errs) {
	var errs Errs
	var rdsTags []*rds.ListTagsForResourceOutput

	for _, svc := range c.svcs {
		if svc.rds == nil {
			svc.rds = rds.New(svc.session)
		}
		rdsTag, err := svc.rds.ListTagsForResource(input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, rds.ServiceName, err))
		} else {
			rdsTags = append(rdsTags, rdsTag)
		}
	}
	return rdsTags, errs
}
