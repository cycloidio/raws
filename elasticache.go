package raws

import (
	"github.com/aws/aws-sdk-go/service/elasticache"
)

// Returns all Elasticache clusters based on the input given
func (c *connector) GetElasticCacheCluster(input *elasticache.DescribeCacheClustersInput) ([]*elasticache.DescribeCacheClustersOutput, Errs) {
	var errs Errs
	var elasticCacheClusters []*elasticache.DescribeCacheClustersOutput

	for _, svc := range c.svcs {
		if svc.elasticache == nil {
			svc.elasticache = elasticache.New(svc.session)
		}
		elasticCacheCluster, err := svc.elasticache.DescribeCacheClusters(input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, elasticache.ServiceName, err))
		} else {
			elasticCacheClusters = append(elasticCacheClusters, elasticCacheCluster)
		}
	}
	return elasticCacheClusters, errs
}

// Returns a list of tags of Elasticache resources based on its ARN
func (c *connector) GetElasticacheTags(input *elasticache.ListTagsForResourceInput) ([]*elasticache.TagListMessage, Errs) {
	var errs Errs
	var elastiCacheTags []*elasticache.TagListMessage

	for _, svc := range c.svcs {
		if svc.elasticache == nil {
			svc.elasticache = elasticache.New(svc.session)
		}
		elasticacheTag, err := svc.elasticache.ListTagsForResource(input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, elasticache.ServiceName, err))
		} else {
			elastiCacheTags = append(elastiCacheTags, elasticacheTag)
		}
	}
	return elastiCacheTags, errs
}
