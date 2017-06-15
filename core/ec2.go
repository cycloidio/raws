package core

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Returns all EC2 instances based on the input given
func (c *Connector) GetInstances(input *ec2.DescribeInstancesInput) ([]*ec2.DescribeInstancesOutput, Errs) {
	var errs Errs
	var instances []*ec2.DescribeInstancesOutput

	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		instance, err := svc.ec2.DescribeInstances(input)
		if err != nil {
			instances = append(instances, instance)
		} else {
			errs = append(errs, NewAPIError(svc.region, ec2.ServiceName, err))
		}
	}
	return instances, errs
}

// Returns all EC2 VPCs based on the input given
func (c *Connector) GetVpcs(input *ec2.DescribeVpcsInput) ([]*ec2.DescribeVpcsOutput, Errs) {
	var errs Errs
	var vpcs []*ec2.DescribeVpcsOutput

	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		vpc, err := svc.ec2.DescribeVpcs(input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, ec2.ServiceName, err))
		} else {
			vpcs = append(vpcs, vpc)
		}
	}
	return vpcs, errs
}

// Returns all EC2 AMI belonging to the Account ID based on the input given
func (c *Connector) GetImages(input *ec2.DescribeImagesInput) ([]*ec2.DescribeImagesOutput, Errs) {
	var errs Errs
	var images []*ec2.DescribeImagesOutput

	if input == nil {
		input = &ec2.DescribeImagesInput{}
	}
	input.Owners = append(input.Owners, c.accountID)
	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		image, err := svc.ec2.DescribeImages(input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, ec2.ServiceName, err))
		} else {
			images = append(images, image)
		}
	}
	return images, errs
}

// Returns all EC2 security groups based on the input given
func (c *Connector) GetSecurityGroups(input *ec2.DescribeSecurityGroupsInput) ([]*ec2.DescribeSecurityGroupsOutput, Errs) {
	var errs Errs
	var secgroups []*ec2.DescribeSecurityGroupsOutput

	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		secgroup, err := svc.ec2.DescribeSecurityGroups(input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, ec2.ServiceName, err))
		} else {
			secgroups = append(secgroups, secgroup)
		}
	}
	return secgroups, errs
}

// Returns all EC2 subnets based on the input given
func (c *Connector) GetSubnets(input *ec2.DescribeSubnetsInput) ([]*ec2.DescribeSubnetsOutput, Errs) {
	var errs Errs
	var subnets []*ec2.DescribeSubnetsOutput

	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		subnet, err := svc.ec2.DescribeSubnets(input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, ec2.ServiceName, err))
		} else {
			subnets = append(subnets, subnet)
		}
	}
	return subnets, errs
}

// Returns all EC2 volumes based on the input given
func (c *Connector) GetVolumes(input *ec2.DescribeVolumesInput) ([]*ec2.DescribeVolumesOutput, Errs) {
	var errs Errs
	var volumes []*ec2.DescribeVolumesOutput

	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		volume, err := svc.ec2.DescribeVolumes(input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, ec2.ServiceName, err))
		} else {
			volumes = append(volumes, volume)
		}
	}
	return volumes, errs
}

// Returns all snapshots belonging to the Account ID based on the input given
func (c *Connector) GetSnapshots(input *ec2.DescribeSnapshotsInput) ([]*ec2.DescribeSnapshotsOutput, Errs) {
	var errs Errs
	var snapshots []*ec2.DescribeSnapshotsOutput

	if input == nil {
		input = &ec2.DescribeSnapshotsInput{}
	}
	input.OwnerIds = append(input.OwnerIds, c.accountID)
	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		snapshot, err := svc.ec2.DescribeSnapshots(input)
		if err != nil {
			snapshots = append(snapshots, snapshot)
		} else {
			errs = append(errs, NewAPIError(svc.region, ec2.ServiceName, err))
		}
	}
	return snapshots, errs
}
