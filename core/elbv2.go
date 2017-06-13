package core

import (
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func (c *Connector) GetLoadBalancersV2(input *elbv2.DescribeLoadBalancersInput) ([]*elbv2.DescribeLoadBalancersOutput, error) {
	var errElbs RawsErr = RawsErr{}
	var elbs []*elbv2.DescribeLoadBalancersOutput

	for _, svc := range c.svcs {
		if svc.elbv2 == nil {
			svc.elbv2 = elbv2.New(svc.session)
		}
		elb, err := svc.elbv2.DescribeLoadBalancers(input)
		elbs = append(elbs, elb)
		errElbs.AppendError(svc.region, elbv2.ServiceName, err)
	}
	return elbs, errElbs
}

func (c *Connector) GetLoadBalancersV2Tags(input *elbv2.DescribeTagsInput) ([]*elbv2.DescribeTagsOutput, error) {
	var errTags RawsErr = RawsErr{}
	var elbTags []*elbv2.DescribeTagsOutput

	for _, svc := range c.svcs {
		if svc.elbv2 == nil {
			svc.elbv2 = elbv2.New(svc.session)
		}
		elbTag, err := svc.elbv2.DescribeTags(input)
		elbTags = append(elbTags, elbTag)
		errTags.AppendError(svc.region, elbv2.ServiceName, err)
	}
	return elbTags, errTags
}
