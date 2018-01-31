package raws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/elb"
)

func (c *connector) GetLoadBalancers(
	ctx context.Context, input *elb.DescribeLoadBalancersInput,
) ([]*elb.DescribeLoadBalancersOutput, Errs) {
	var elbs []*elb.DescribeLoadBalancersOutput
	var errs Errs

	for _, svc := range c.svcs {
		if svc.elb == nil {
			svc.elb = elb.New(svc.session)
		}
		elbv1, err := svc.elb.DescribeLoadBalancersWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, elb.ServiceName, err))
		} else {
			elbs = append(elbs, elbv1)
		}
	}
	return elbs, errs
}

func (c *connector) GetLoadBalancersTags(
	ctx context.Context, input *elb.DescribeTagsInput,
) ([]*elb.DescribeTagsOutput, Errs) {
	var elbTags []*elb.DescribeTagsOutput
	var errs Errs

	for _, svc := range c.svcs {
		if svc.elb == nil {
			svc.elb = elb.New(svc.session)
		}
		tags, err := svc.elb.DescribeTagsWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, elb.ServiceName, err))
		} else {
			elbTags = append(elbTags, tags)
		}
	}
	return elbTags, errs
}
