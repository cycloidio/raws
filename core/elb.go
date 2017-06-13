package core

import (
	"github.com/aws/aws-sdk-go/service/elb"
)

// Returns a list of ELB (v1) based on the input from the different regions
func (c *Connector) GetLoadBalancers(input *elb.DescribeLoadBalancersInput) ([]*elb.DescribeLoadBalancersOutput, error) {
	var elbs []*elb.DescribeLoadBalancersOutput
	var errElbs RawsErr = RawsErr{}

	for _, svc := range c.svcs {
		if svc.elb == nil {
			svc.elb = elb.New(svc.session)
		}
		elbv1, err := svc.elb.DescribeLoadBalancers(input)
		elbs = append(elbs, elbv1)
		errElbs.AppendError(svc.region, elb.ServiceName, err)
	}
	if len(errElbs.APIErrs) == 0 {
		return elbs, nil
	}
	return elbs, errElbs
}

// Returns a list of Tags based on the input from the different regions
func (c *Connector) GetLoadBalancersTags(input *elb.DescribeTagsInput) ([]*elb.DescribeTagsOutput, error) {
	var elbTags []*elb.DescribeTagsOutput
	var errTags RawsErr = RawsErr{}

	for _, svc := range c.svcs {
		if svc.elb == nil {
			svc.elb = elb.New(svc.session)
		}
		tags, err := svc.elb.DescribeTags(input)
		elbTags = append(elbTags, tags)
		errTags.AppendError(svc.region, elb.ServiceName, err)
	}
	if len(errTags.APIErrs) == 0 {
		return elbTags, nil
	}
	return elbTags, errTags
}
