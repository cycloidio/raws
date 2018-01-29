package raws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/elbv2"
)

func (c *connector) GetLoadBalancersV2(
	ctx context.Context, input *elbv2.DescribeLoadBalancersInput,
) ([]*elbv2.DescribeLoadBalancersOutput, Errs) {
	var errs Errs
	var elbs []*elbv2.DescribeLoadBalancersOutput

	for _, svc := range c.svcs {
		if svc.elbv2 == nil {
			svc.elbv2 = elbv2.New(svc.session)
		}
		elb, err := svc.elbv2.DescribeLoadBalancersWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, elbv2.ServiceName, err))
		} else {
			elbs = append(elbs, elb)
		}
	}
	return elbs, errs
}

func (c *connector) GetLoadBalancersV2Tags(
	ctx context.Context, input *elbv2.DescribeTagsInput,
) ([]*elbv2.DescribeTagsOutput, Errs) {
	var errs Errs
	var elbTags []*elbv2.DescribeTagsOutput

	for _, svc := range c.svcs {
		if svc.elbv2 == nil {
			svc.elbv2 = elbv2.New(svc.session)
		}
		elbTag, err := svc.elbv2.DescribeTagsWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, elbv2.ServiceName, err))
		} else {
			elbTags = append(elbTags, elbTag)
		}
	}
	return elbTags, errs
}
