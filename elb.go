package raws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/elb"
)

func (c *connector) GetLoadBalancers(
	ctx context.Context, input *elb.DescribeLoadBalancersInput,
) (map[string]elb.DescribeLoadBalancersOutput, error) {
	var errs Errors
	var elbs = map[string]elb.DescribeLoadBalancersOutput{}

	for _, svc := range c.svcs {
		if svc.elb == nil {
			svc.elb = elb.New(svc.session)
		}
		elbv1, err := svc.elb.DescribeLoadBalancersWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewError(svc.region, elb.ServiceName, err))
		} else {
			elbs[svc.region] = *elbv1
		}
	}

	if errs != nil {
		return elbs, errs
	}

	return elbs, nil
}

func (c *connector) GetLoadBalancersTags(
	ctx context.Context, input *elb.DescribeTagsInput,
) ([]*elb.DescribeTagsOutput, error) {
	var elbTags []*elb.DescribeTagsOutput
	var errs Errors

	for _, svc := range c.svcs {
		if svc.elb == nil {
			svc.elb = elb.New(svc.session)
		}
		tags, err := svc.elb.DescribeTagsWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewError(svc.region, elb.ServiceName, err))
		} else {
			elbTags = append(elbTags, tags)
		}
	}

	if errs != nil {
		return elbTags, errs
	}

	return elbTags, nil
}
