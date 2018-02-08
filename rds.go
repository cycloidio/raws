package raws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/rds"
)

func (c *connector) GetDBInstances(
	ctx context.Context, input *rds.DescribeDBInstancesInput,
) (map[string]rds.DescribeDBInstancesOutput, error) {
	var errs Errors
	var instances = map[string]rds.DescribeDBInstancesOutput{}

	for _, svc := range c.svcs {
		if svc.rds == nil {
			svc.rds = rds.New(svc.session)
		}
		instance, err := svc.rds.DescribeDBInstancesWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewError(svc.region, rds.ServiceName, err))
		} else {
			instances[svc.region] = *instance
		}
	}

	if errs != nil {
		return instances, errs
	}

	return instances, nil
}

func (c *connector) GetDBInstancesTags(
	ctx context.Context, input *rds.ListTagsForResourceInput,
) (map[string]rds.ListTagsForResourceOutput, error) {
	var errs Errors
	var rdsTags = map[string]rds.ListTagsForResourceOutput{}

	for _, svc := range c.svcs {
		if svc.rds == nil {
			svc.rds = rds.New(svc.session)
		}
		rdsTag, err := svc.rds.ListTagsForResourceWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewError(svc.region, rds.ServiceName, err))
		} else {
			rdsTags[svc.region] = *rdsTag
		}
	}

	if errs != nil {
		return rdsTags, errs
	}

	return rdsTags, nil
}
