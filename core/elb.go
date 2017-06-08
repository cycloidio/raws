package core

import (
	"github.com/aws/aws-sdk-go/service/elb"
)

func (c *Connector) GetLoadBalancers(input *elb.DescribeLoadBalancersInput) ([]*elb.DescribeLoadBalancersOutput, []error) {
	var errs []error
	var elbs []*elb.DescribeLoadBalancersOutput

	for _, svc := range c.svcs {
		if svc.elb == nil {
			svc.elb = elb.New(svc.session)
		}
		elb, err := svc.elb.DescribeLoadBalancers(input)
		elbs = append(elbs, elb)
		errs = append(errs, err)
	}
	return elbs, errs
}

func (c *Connector) GetLoadBalancersTags(input *elb.DescribeTagsInput) ([]*elb.DescribeTagsOutput, []error) {
	var errs []error
	var elbTags []*elb.DescribeTagsOutput

	for _, svc := range c.svcs {
		if svc.elb == nil {
			svc.elb = elb.New(svc.session)
		}
		elb, err := svc.elb.DescribeTags(input)
		elbTags = append(elbTags, elb)
		errs = append(errs, err)
	}
	return elbTags, errs
}
