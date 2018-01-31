package raws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/elasticache"
)

func (c *connector) GetElastiCacheCluster(
	ctx context.Context, input *elasticache.DescribeCacheClustersInput,
) ([]*elasticache.DescribeCacheClustersOutput, Errs) {
	var errs Errs
	var elasticCacheClusters []*elasticache.DescribeCacheClustersOutput

	for _, svc := range c.svcs {
		if svc.elasticache == nil {
			svc.elasticache = elasticache.New(svc.session)
		}
		elasticCacheCluster, err := svc.elasticache.DescribeCacheClustersWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, elasticache.ServiceName, err))
		} else {
			elasticCacheClusters = append(elasticCacheClusters, elasticCacheCluster)
		}
	}
	return elasticCacheClusters, errs
}

func (c *connector) GetElastiCacheTags(
	ctx context.Context, input *elasticache.ListTagsForResourceInput,
) ([]*elasticache.TagListMessage, Errs) {
	var errs Errs
	var elastiCacheTags []*elasticache.TagListMessage

	for _, svc := range c.svcs {
		if svc.elasticache == nil {
			svc.elasticache = elasticache.New(svc.session)
		}
		elasticacheTag, err := svc.elasticache.ListTagsForResourceWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, elasticache.ServiceName, err))
		} else {
			elastiCacheTags = append(elastiCacheTags, elasticacheTag)
		}
	}
	return elastiCacheTags, errs
}
