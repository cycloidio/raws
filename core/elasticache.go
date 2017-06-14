package core

import (
	"github.com/aws/aws-sdk-go/service/elasticache"
)

// Returns all Elasticache clusters based on the input given
func (c *Connector) GetElasticCacheCluster(input *elasticache.DescribeCacheClustersInput) ([]*elasticache.DescribeCacheClustersOutput, error) {
	var errs = RawsErr{}
	var elasticCacheClusters []*elasticache.DescribeCacheClustersOutput

	for _, svc := range c.svcs {
		if svc.elasticache == nil {
			svc.elasticache = elasticache.New(svc.session)
		}
		elasticCacheCluster, err := svc.elasticache.DescribeCacheClusters(input)
		elasticCacheClusters = append(elasticCacheClusters, elasticCacheCluster)
		errs.AppendError(svc.region, elasticache.ServiceName, err)
	}
	if len(errs.APIErrs) == 0 {
		return elasticCacheClusters, nil
	}
	return elasticCacheClusters, errs
}

// Returns a list of tags of Elasticache resources based on its ARN
func (c *Connector) GetElasticacheTags(input *elasticache.ListTagsForResourceInput) ([]*elasticache.TagListMessage, error) {
	var errs RawsErr = RawsErr{}
	var elastiCacheTags []*elasticache.TagListMessage

	for _, svc := range c.svcs {
		if svc.elasticache == nil {
			svc.elasticache = elasticache.New(svc.session)
		}
		elasticacheTag, err := svc.elasticache.ListTagsForResource(input)
		elastiCacheTags = append(elastiCacheTags, elasticacheTag)
		errs.AppendError(svc.region, elasticache.ServiceName, err)
	}
	if len(errs.APIErrs) == 0 {
		return elastiCacheTags, nil
	}
	return elastiCacheTags, errs
}
