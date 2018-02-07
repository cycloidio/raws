package raws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/rds"
)

func (c *connector) GetDBInstances(
	ctx context.Context, input *rds.DescribeDBInstancesInput,
) ([]*rds.DescribeDBInstancesOutput, error) {
	var errs Errs
	var instances []*rds.DescribeDBInstancesOutput

	for _, svc := range c.svcs {
		if svc.rds == nil {
			svc.rds = rds.New(svc.session)
		}
		instance, err := svc.rds.DescribeDBInstancesWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, rds.ServiceName, err))
		} else {
			instances = append(instances, instance)
		}
	}

	if errs != nil {
		return instances, errs
	}

	return instances, nil
}

func (c *connector) GetDBInstancesTags(
	ctx context.Context, input *rds.ListTagsForResourceInput,
) ([]*rds.ListTagsForResourceOutput, Errs) {
	var errs Errs
	var rdsTags []*rds.ListTagsForResourceOutput

	for _, svc := range c.svcs {
		if svc.rds == nil {
			svc.rds = rds.New(svc.session)
		}
		rdsTag, err := svc.rds.ListTagsForResourceWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, rds.ServiceName, err))
		} else {
			rdsTags = append(rdsTags, rdsTag)
		}
	}
	return rdsTags, errs
}
