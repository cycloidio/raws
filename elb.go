package raws

import (
	"github.com/aws/aws-sdk-go/service/elb"
)

// Returns a list of ELB (v1) based on the input from the different regions
func (c *connector) GetLoadBalancers(input *elb.DescribeLoadBalancersInput) ([]*elb.DescribeLoadBalancersOutput, Errs) {
	var elbs []*elb.DescribeLoadBalancersOutput
	var errs Errs

	for _, svc := range c.svcs {
		if svc.elb == nil {
			svc.elb = elb.New(svc.session)
		}
		elbv1, err := svc.elb.DescribeLoadBalancers(input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, elb.ServiceName, err))
		} else {
			elbs = append(elbs, elbv1)
		}
	}
	return elbs, errs
}

// Returns a list of Tags based on the input from the different regions
func (c *connector) GetLoadBalancersTags(input *elb.DescribeTagsInput) ([]*elb.DescribeTagsOutput, Errs) {
	var elbTags []*elb.DescribeTagsOutput
	var errs Errs

	for _, svc := range c.svcs {
		if svc.elb == nil {
			svc.elb = elb.New(svc.session)
		}
		tags, err := svc.elb.DescribeTags(input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, elb.ServiceName, err))
		} else {
			elbTags = append(elbTags, tags)
		}
	}
	return elbTags, errs
}
