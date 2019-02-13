package raws

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/elasticache/elasticacheiface"
)

type mockElasticCache struct {
	elasticacheiface.ElastiCacheAPI

	// Mocking of DescribeCacheClusters
	dcco   *elasticache.DescribeCacheClustersOutput
	dccerr error

	// Mocking of ListTagsForResource
	ltfro   *elasticache.TagListMessage
	ltfrerr error
}

func (m mockElasticCache) DescribeCacheClustersWithContext(
	_ aws.Context, _ *elasticache.DescribeCacheClustersInput, _ ...request.Option,
) (*elasticache.DescribeCacheClustersOutput, error) {
	return m.dcco, m.dccerr
}

func (m mockElasticCache) ListTagsForResourceWithContext(
	_ aws.Context, _ *elasticache.ListTagsForResourceInput, _ ...request.Option,
) (*elasticache.TagListMessage, error) {
	return m.ltfro, m.ltfrerr
}

func TestGetElastiCacheCluster(t *testing.T) {
	tests := []struct {
		name             string
		mocked           []*serviceConnector
		expectedClusters map[string]elasticache.DescribeCacheClustersOutput
		expectedError    error
	}{{name: "one region no error",
		mocked: []*serviceConnector{
			{
				region: "test",
				elasticache: mockElasticCache{
					dcco: &elasticache.DescribeCacheClustersOutput{
						CacheClusters: []*elasticache.CacheCluster{
							{
								CacheClusterId: aws.String("1"),
							},
						},
					},
					dccerr: nil,
				},
			},
		},
		expectedError: nil,
		expectedClusters: map[string]elasticache.DescribeCacheClustersOutput{
			"test": {
				CacheClusters: []*elasticache.CacheCluster{
					{
						CacheClusterId: aws.String("1"),
					},
				},
			},
		},
	},
		{name: "one region with error",
			mocked: []*serviceConnector{
				{
					region: "test",
					elasticache: mockElasticCache{
						dcco:   &elasticache.DescribeCacheClustersOutput{},
						dccerr: errors.New("error with test"),
					},
				},
			},
			expectedError: Errors{Error{
				err:     errors.New("error with test"),
				region:  "test",
				service: elasticache.ServiceName,
			}},
			expectedClusters: map[string]elasticache.DescribeCacheClustersOutput{},
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					elasticache: mockElasticCache{
						dcco: &elasticache.DescribeCacheClustersOutput{
							CacheClusters: []*elasticache.CacheCluster{
								{
									CacheClusterId: aws.String("1"),
								},
							},
						},
						dccerr: nil,
					},
				},
				{
					region: "test-2",
					elasticache: mockElasticCache{
						dcco: &elasticache.DescribeCacheClustersOutput{
							CacheClusters: []*elasticache.CacheCluster{
								{
									CacheClusterId: aws.String("2"),
								},
							},
						},
						dccerr: nil,
					},
				},
			},
			expectedError: nil,
			expectedClusters: map[string]elasticache.DescribeCacheClustersOutput{
				"test-1": {
					CacheClusters: []*elasticache.CacheCluster{
						{
							CacheClusterId: aws.String("1"),
						},
					},
				},
				"test-2": {
					CacheClusters: []*elasticache.CacheCluster{
						{
							CacheClusterId: aws.String("2"),
						},
					},
				},
			},
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					elasticache: mockElasticCache{
						dcco:   &elasticache.DescribeCacheClustersOutput{},
						dccerr: errors.New("error with test"),
					},
				},
				{
					region: "test-2",
					elasticache: mockElasticCache{
						dcco: &elasticache.DescribeCacheClustersOutput{
							CacheClusters: []*elasticache.CacheCluster{
								{
									CacheClusterId: aws.String("2"),
								},
							},
						},
						dccerr: nil,
					},
				},
			},
			expectedError: Errors{
				Error{
					err:     errors.New("error with test"),
					region:  "test-1",
					service: elasticache.ServiceName,
				},
			},
			expectedClusters: map[string]elasticache.DescribeCacheClustersOutput{
				"test-2": {
					CacheClusters: []*elasticache.CacheCluster{
						{
							CacheClusterId: aws.String("2"),
						},
					},
				},
			},
		},
	}

	var ctx = context.Background()

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		cluster, err := c.GetElastiCacheClusters(ctx, nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(cluster, tt.expectedClusters) {
			t.Errorf("%s [%d] - clusters: received=%+v | expected=%+v",
				tt.name, i, cluster, tt.expectedClusters)
		}
	}
}

func TestGetElastiCacheTags(t *testing.T) {
	tests := []struct {
		name          string
		mocked        []*serviceConnector
		expectedTags  map[string]elasticache.TagListMessage
		expectedError error
	}{{name: "one region no error",
		mocked: []*serviceConnector{
			{
				region: "test",
				elasticache: mockElasticCache{
					ltfro: &elasticache.TagListMessage{
						TagList: []*elasticache.Tag{
							{
								Key:   aws.String("test"),
								Value: aws.String("1"),
							}},
					},
					ltfrerr: nil,
				},
			},
		},
		expectedError: nil,
		expectedTags: map[string]elasticache.TagListMessage{
			"test": {
				TagList: []*elasticache.Tag{
					{
						Key:   aws.String("test"),
						Value: aws.String("1"),
					},
				},
			},
		},
	},
		{name: "one region with error",
			mocked: []*serviceConnector{
				{
					region: "test",
					elasticache: mockElasticCache{
						ltfro: &elasticache.TagListMessage{
							TagList: []*elasticache.Tag{},
						},
						ltfrerr: errors.New("error with test"),
					},
				},
			},
			expectedError: Errors{
				Error{
					err:     errors.New("error with test"),
					region:  "test",
					service: elasticache.ServiceName,
				},
			},
			expectedTags: map[string]elasticache.TagListMessage{},
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					elasticache: mockElasticCache{
						ltfro: &elasticache.TagListMessage{
							TagList: []*elasticache.Tag{
								{
									Key:   aws.String("test"),
									Value: aws.String("1"),
								}},
						},
						ltfrerr: nil,
					},
				},
				{
					region: "test-2",
					elasticache: mockElasticCache{
						ltfro: &elasticache.TagListMessage{
							TagList: []*elasticache.Tag{
								{
									Key:   aws.String("test"),
									Value: aws.String("2"),
								}},
						},
						ltfrerr: nil,
					},
				},
			},
			expectedError: nil,
			expectedTags: map[string]elasticache.TagListMessage{
				"test-1": {
					TagList: []*elasticache.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("1"),
						},
					},
				},
				"test-2": {
					TagList: []*elasticache.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("2"),
						},
					},
				},
			},
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					elasticache: mockElasticCache{
						ltfro:   &elasticache.TagListMessage{},
						ltfrerr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					elasticache: mockElasticCache{
						ltfro: &elasticache.TagListMessage{
							TagList: []*elasticache.Tag{
								{
									Key:   aws.String("test"),
									Value: aws.String("2"),
								}},
						},
						ltfrerr: nil,
					},
				},
			},
			expectedError: Errors{
				Error{
					err:     errors.New("error with test-1"),
					region:  "test-1",
					service: elasticache.ServiceName,
				},
			},
			expectedTags: map[string]elasticache.TagListMessage{
				"test-2": {
					TagList: []*elasticache.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("2"),
						},
					},
				},
			},
		},
	}

	var ctx = context.Background()

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		tags, err := c.GetElastiCacheTags(ctx, nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(tags, tt.expectedTags) {
			t.Errorf("%s [%d] - tags: received=%+v | expected=%+v",
				tt.name, i, tags, tt.expectedTags)
		}
	}
}
