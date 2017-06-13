package core

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (c *Connector) GetInstances(input *ec2.DescribeInstancesInput) ([]*ec2.DescribeInstancesOutput, error) {
	var errInsts RawsErr = RawsErr{}
	var instances []*ec2.DescribeInstancesOutput

	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		instance, err := svc.ec2.DescribeInstances(input)
		instances = append(instances, instance)
		errInsts.AppendError(svc.region, ec2.ServiceName, err)
	}
	return instances, errInsts
}

func (c *Connector) GetVpcs(input *ec2.DescribeVpcsInput) ([]*ec2.DescribeVpcsOutput, error) {
	var errInsts RawsErr = RawsErr{}
	var instances []*ec2.DescribeVpcsOutput

	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		instance, err := svc.ec2.DescribeVpcs(input)
		instances = append(instances, instance)
		errInsts.AppendError(svc.region, ec2.ServiceName, err)
	}
	return instances, errInsts
}

func (c *Connector) GetImages(input *ec2.DescribeImagesInput) ([]*ec2.DescribeImagesOutput, error) {
	var errInsts RawsErr = RawsErr{}
	var instances []*ec2.DescribeImagesOutput

	if input == nil {
		input = &ec2.DescribeImagesInput{}
	}
	input.Owners = append(input.Owners, c.accountID)
	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		instance, err := svc.ec2.DescribeImages(input)
		instances = append(instances, instance)
		errInsts.AppendError(svc.region, ec2.ServiceName, err)
	}
	return instances, errInsts
}

func (c *Connector) GetSecurityGroups(input *ec2.DescribeSecurityGroupsInput) ([]*ec2.DescribeSecurityGroupsOutput, error) {
	var errInsts RawsErr = RawsErr{}
	var instances []*ec2.DescribeSecurityGroupsOutput

	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		instance, err := svc.ec2.DescribeSecurityGroups(input)
		instances = append(instances, instance)
		errInsts.AppendError(svc.region, ec2.ServiceName, err)
	}
	return instances, errInsts
}

func (c *Connector) GetSubnets(input *ec2.DescribeSubnetsInput) ([]*ec2.DescribeSubnetsOutput, error) {
	var errInsts RawsErr = RawsErr{}
	var instances []*ec2.DescribeSubnetsOutput

	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		instance, err := svc.ec2.DescribeSubnets(input)
		instances = append(instances, instance)
		errInsts.AppendError(svc.region, ec2.ServiceName, err)
	}
	return instances, errInsts
}

func (c *Connector) GetVolumes(input *ec2.DescribeVolumesInput) ([]*ec2.DescribeVolumesOutput, error) {
	var errInsts RawsErr = RawsErr{}
	var instances []*ec2.DescribeVolumesOutput

	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		instance, err := svc.ec2.DescribeVolumes(input)
		instances = append(instances, instance)
		errInsts.AppendError(svc.region, ec2.ServiceName, err)
	}
	return instances, errInsts
}

func (c *Connector) GetSnapshots(input *ec2.DescribeSnapshotsInput) ([]*ec2.DescribeSnapshotsOutput, error) {
	var errInsts RawsErr = RawsErr{}
	var instances []*ec2.DescribeSnapshotsOutput

	if input == nil {
		input = &ec2.DescribeSnapshotsInput{}
	}
	input.OwnerIds = append(input.OwnerIds, c.accountID)
	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		instance, err := svc.ec2.DescribeSnapshots(input)
		instances = append(instances, instance)
		errInsts.AppendError(svc.region, ec2.ServiceName, err)
	}
	return instances, errInsts
}
