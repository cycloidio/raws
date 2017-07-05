package raws

import (
	"github.com/aws/aws-sdk-go/service/elbv2"
)

// Returns a list of ELB (v2) - also known as ALB - based on the input from the different regions
func (c *Connector) GetLoadBalancersV2(input *elbv2.DescribeLoadBalancersInput) ([]*elbv2.DescribeLoadBalancersOutput, Errs) {
	var errs Errs
	var elbs []*elbv2.DescribeLoadBalancersOutput

	for _, svc := range c.svcs {
		if svc.elbv2 == nil {
			svc.elbv2 = elbv2.New(svc.session)
		}
		elb, err := svc.elbv2.DescribeLoadBalancers(input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, elbv2.ServiceName, err))
		} else {
			elbs = append(elbs, elb)
		}
	}
	return elbs, errs
}

// Returns a list of Tags based on the input from the different regions
func (c *Connector) GetLoadBalancersV2Tags(input *elbv2.DescribeTagsInput) ([]*elbv2.DescribeTagsOutput, Errs) {
	var errs Errs
	var elbTags []*elbv2.DescribeTagsOutput

	for _, svc := range c.svcs {
		if svc.elbv2 == nil {
			svc.elbv2 = elbv2.New(svc.session)
		}
		elbTag, err := svc.elbv2.DescribeTags(input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, elbv2.ServiceName, err))
		} else {
			elbTags = append(elbTags, elbTag)
		}
	}
	return elbTags, errs
}
