package raws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/configservice"
)

func (c *connector) GetRecordedResourceCounts(
	ctx context.Context, input *configservice.GetDiscoveredResourceCountsInput,
) (map[string]configservice.GetDiscoveredResourceCountsOutput, error) {
	var errs Errors
	var resCounts = map[string]configservice.GetDiscoveredResourceCountsOutput{}

	for _, svc := range c.svcs {
		if svc.configservice == nil {
			svc.configservice = configservice.New(svc.session)
		}
		counts, err := svc.configservice.GetDiscoveredResourceCountsWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewError(svc.region, configservice.ServiceName, err))
		} else {
			resCounts[svc.region] = *counts
		}
	}

	if errs != nil {
		return resCounts, errs
	}

	return resCounts, nil
}
