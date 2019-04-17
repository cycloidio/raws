package raws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/cloudfront"
)

func (c *connector) GetCloudFrontDistributions(
	ctx context.Context, input *cloudfront.ListDistributionsInput,
) (map[string]cloudfront.ListDistributionsOutput, error) {
	var errs Errors
	var regionDistributions = map[string]cloudfront.ListDistributionsOutput{}

	for _, svc := range c.svcs {
		if svc.cloudfront == nil {
			svc.cloudfront = cloudfront.New(svc.session)
		}
		distributions, err := svc.cloudfront.ListDistributionsWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewError(svc.region, cloudfront.ServiceName, err))
		} else {
			regionDistributions[svc.region] = *distributions
		}
	}

	if errs != nil {
		return regionDistributions, errs
	}

	return regionDistributions, nil
}

func (c *connector) GetCloudFrontPublicKeys(
	ctx context.Context, input *cloudfront.ListPublicKeysInput,
) (map[string]cloudfront.ListPublicKeysOutput, error) {
	var errs Errors
	var regionPublicKeys = map[string]cloudfront.ListPublicKeysOutput{}

	for _, svc := range c.svcs {
		if svc.cloudfront == nil {
			svc.cloudfront = cloudfront.New(svc.session)
		}
		publicKeys, err := svc.cloudfront.ListPublicKeysWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewError(svc.region, cloudfront.ServiceName, err))
		} else {
			regionPublicKeys[svc.region] = *publicKeys
		}
	}

	if errs != nil {
		return regionPublicKeys, errs
	}

	return regionPublicKeys, nil
}
