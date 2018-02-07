package raws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func (c *connector) GetInstances(
	ctx context.Context, input *ec2.DescribeInstancesInput,
) (map[string]ec2.DescribeInstancesOutput, error) {
	var errs Errors
	var instances = map[string]ec2.DescribeInstancesOutput{}

	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		instance, err := svc.ec2.DescribeInstancesWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewError(svc.region, ec2.ServiceName, err))
		} else {
			instances[svc.region] = *instance
		}
	}

	if errs != nil {
		return instances, errs
	}

	return instances, nil
}

func (c *connector) GetVpcs(
	ctx context.Context, input *ec2.DescribeVpcsInput,
) (map[string]ec2.DescribeVpcsOutput, error) {
	var errs Errors
	var vpcs = map[string]ec2.DescribeVpcsOutput{}

	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		vpc, err := svc.ec2.DescribeVpcsWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewError(svc.region, ec2.ServiceName, err))
		} else {
			vpcs[svc.region] = *vpc
		}
	}

	if errs != nil {
		return vpcs, errs
	}

	return vpcs, nil
}

func (c *connector) GetImages(
	ctx context.Context, input *ec2.DescribeImagesInput,
) (map[string]ec2.DescribeImagesOutput, error) {
	var errs Errors
	var images = map[string]ec2.DescribeImagesOutput{}

	if input == nil {
		input = &ec2.DescribeImagesInput{}
	}
	input.Owners = append(input.Owners, c.accountID)
	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		image, err := svc.ec2.DescribeImagesWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewError(svc.region, ec2.ServiceName, err))
		} else {
			images[svc.region] = *image
		}
	}

	if errs != nil {
		return images, errs
	}

	return images, nil
}

func (c *connector) GetSecurityGroups(
	ctx context.Context, input *ec2.DescribeSecurityGroupsInput,
) (map[string]ec2.DescribeSecurityGroupsOutput, error) {
	var errs Errors
	var secgroups = map[string]ec2.DescribeSecurityGroupsOutput{}

	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		secgroup, err := svc.ec2.DescribeSecurityGroupsWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewError(svc.region, ec2.ServiceName, err))
		} else {
			secgroups[svc.region] = *secgroup
		}
	}

	if errs != nil {
		return secgroups, errs
	}

	return secgroups, nil
}

func (c *connector) GetSubnets(
	ctx context.Context, input *ec2.DescribeSubnetsInput,
) (map[string]ec2.DescribeSubnetsOutput, error) {
	var errs Errors
	var subnets = map[string]ec2.DescribeSubnetsOutput{}

	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		subnet, err := svc.ec2.DescribeSubnetsWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewError(svc.region, ec2.ServiceName, err))
		} else {
			subnets[svc.region] = *subnet
		}
	}

	if errs != nil {
		return subnets, errs
	}

	return subnets, nil
}

func (c *connector) GetVolumes(
	ctx context.Context, input *ec2.DescribeVolumesInput,
) ([]*ec2.DescribeVolumesOutput, error) {
	var errs Errors
	var volumes []*ec2.DescribeVolumesOutput

	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		volume, err := svc.ec2.DescribeVolumesWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewError(svc.region, ec2.ServiceName, err))
		} else {
			volumes = append(volumes, volume)
		}
	}

	if errs != nil {
		return volumes, errs
	}

	return volumes, nil
}

func (c *connector) GetSnapshots(
	ctx context.Context, input *ec2.DescribeSnapshotsInput,
) ([]*ec2.DescribeSnapshotsOutput, error) {
	var errs Errors
	var snapshots []*ec2.DescribeSnapshotsOutput

	if input == nil {
		input = &ec2.DescribeSnapshotsInput{}
	}
	input.OwnerIds = append(input.OwnerIds, c.accountID)
	for _, svc := range c.svcs {
		if svc.ec2 == nil {
			svc.ec2 = ec2.New(svc.session)
		}
		snapshot, err := svc.ec2.DescribeSnapshotsWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewError(svc.region, ec2.ServiceName, err))
		} else {
			snapshots = append(snapshots, snapshot)
		}
	}

	if errs != nil {
		return snapshots, errs
	}

	return snapshots, nil
}
