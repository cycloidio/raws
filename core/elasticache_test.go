package core

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
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

func (c mockElasticCache) DescribeCacheClusters(input *elasticache.DescribeCacheClustersInput) (*elasticache.DescribeCacheClustersOutput, error) {
	return c.dcco, c.dccerr
}

func (c mockElasticCache) ListTagsForResource(input *elasticache.ListTagsForResourceInput) (*elasticache.TagListMessage, error) {
	return c.ltfro, c.ltfrerr
}

func TestGetElasticacheCluster(t *testing.T) {
	tests := []struct {
		name             string
		mocked           []*serviceConnector
		expectedClusters []*elasticache.DescribeCacheClustersOutput
		expectedError    Errs
	}{{name: "one region no error",
		mocked: []*serviceConnector{
			&serviceConnector{
				region: "test",
				elasticache: mockElasticCache{
					dcco: &elasticache.DescribeCacheClustersOutput{
						CacheClusters: []*elasticache.CacheCluster{
							&elasticache.CacheCluster{
								CacheClusterId: aws.String("1"),
							},
						},
					},
					dccerr: nil,
				},
			},
		},
		expectedError: nil,
		expectedClusters: []*elasticache.DescribeCacheClustersOutput{
			&elasticache.DescribeCacheClustersOutput{
				CacheClusters: []*elasticache.CacheCluster{
					&elasticache.CacheCluster{
						CacheClusterId: aws.String("1"),
					},
				},
			},
		},
	},
		{name: "one region with error",
			mocked: []*serviceConnector{
				&serviceConnector{
					region: "test",
					elasticache: mockElasticCache{
						dcco:   &elasticache.DescribeCacheClustersOutput{},
						dccerr: errors.New("error with test"),
					},
				},
			},
			expectedError: Errs{&callErr{
				err:     errors.New("error with test"),
				region:  "test",
				service: elasticache.ServiceName,
			}},
			expectedClusters: nil,
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				&serviceConnector{
					region: "test-1",
					elasticache: mockElasticCache{
						dcco: &elasticache.DescribeCacheClustersOutput{
							CacheClusters: []*elasticache.CacheCluster{
								&elasticache.CacheCluster{
									CacheClusterId: aws.String("1"),
								},
							},
						},
						dccerr: nil,
					},
				},
				&serviceConnector{
					region: "test-2",
					elasticache: mockElasticCache{
						dcco: &elasticache.DescribeCacheClustersOutput{
							CacheClusters: []*elasticache.CacheCluster{
								&elasticache.CacheCluster{
									CacheClusterId: aws.String("2"),
								},
							},
						},
						dccerr: nil,
					},
				},
			},
			expectedError: nil,
			expectedClusters: []*elasticache.DescribeCacheClustersOutput{
				&elasticache.DescribeCacheClustersOutput{
					CacheClusters: []*elasticache.CacheCluster{
						&elasticache.CacheCluster{
							CacheClusterId: aws.String("1"),
						},
					},
				},
				&elasticache.DescribeCacheClustersOutput{
					CacheClusters: []*elasticache.CacheCluster{
						&elasticache.CacheCluster{
							CacheClusterId: aws.String("2"),
						},
					},
				},
			},
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				&serviceConnector{
					region: "test-1",
					elasticache: mockElasticCache{
						dcco:   &elasticache.DescribeCacheClustersOutput{},
						dccerr: errors.New("error with test"),
					},
				},
				&serviceConnector{
					region: "test-2",
					elasticache: mockElasticCache{
						dcco: &elasticache.DescribeCacheClustersOutput{
							CacheClusters: []*elasticache.CacheCluster{
								&elasticache.CacheCluster{
									CacheClusterId: aws.String("2"),
								},
							},
						},
						dccerr: nil,
					},
				},
			},
			expectedError: Errs{
				&callErr{
					err:     errors.New("error with test"),
					region:  "test-1",
					service: elasticache.ServiceName,
				},
			},
			expectedClusters: []*elasticache.DescribeCacheClustersOutput{
				&elasticache.DescribeCacheClustersOutput{
					CacheClusters: []*elasticache.CacheCluster{
						&elasticache.CacheCluster{
							CacheClusterId: aws.String("2"),
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		c := &Connector{svcs: tt.mocked}
		cluster, err := c.GetElasticCacheCluster(nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(cluster, tt.expectedClusters) {
			t.Errorf("%s [%d] - clusters: received=%+v | expected=%+v",
				tt.name, i, cluster, tt.expectedClusters)
		}
	}
}

func TestGetElasticacheTags(t *testing.T) {
	tests := []struct {
		name          string
		mocked        []*serviceConnector
		expectedTags  []*elasticache.TagListMessage
		expectedError Errs
	}{{name: "one region no error",
		mocked: []*serviceConnector{
			&serviceConnector{
				region: "test",
				elasticache: mockElasticCache{
					ltfro: &elasticache.TagListMessage{
						TagList: []*elasticache.Tag{
							&elasticache.Tag{
								Key:   aws.String("test"),
								Value: aws.String("1"),
							}},
					},
					ltfrerr: nil,
				},
			},
		},
		expectedError: nil,
		expectedTags: []*elasticache.TagListMessage{
			&elasticache.TagListMessage{
				TagList: []*elasticache.Tag{
					&elasticache.Tag{
						Key:   aws.String("test"),
						Value: aws.String("1"),
					},
				},
			},
		},
	},
		{name: "one region with error",
			mocked: []*serviceConnector{
				&serviceConnector{
					region: "test",
					elasticache: mockElasticCache{
						ltfro: &elasticache.TagListMessage{
							TagList: []*elasticache.Tag{},
						},
						ltfrerr: errors.New("error with test"),
					},
				},
			},
			expectedError: Errs{
				&callErr{
					err:     errors.New("error with test"),
					region:  "test",
					service: elasticache.ServiceName,
				},
			},
			expectedTags: nil,
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				&serviceConnector{
					region: "test-1",
					elasticache: mockElasticCache{
						ltfro: &elasticache.TagListMessage{
							TagList: []*elasticache.Tag{
								&elasticache.Tag{
									Key:   aws.String("test"),
									Value: aws.String("1"),
								}},
						},
						ltfrerr: nil,
					},
				},
				&serviceConnector{
					region: "test-2",
					elasticache: mockElasticCache{
						ltfro: &elasticache.TagListMessage{
							TagList: []*elasticache.Tag{
								&elasticache.Tag{
									Key:   aws.String("test"),
									Value: aws.String("2"),
								}},
						},
						ltfrerr: nil,
					},
				},
			},
			expectedError: nil,
			expectedTags: []*elasticache.TagListMessage{
				&elasticache.TagListMessage{
					TagList: []*elasticache.Tag{
						&elasticache.Tag{
							Key:   aws.String("test"),
							Value: aws.String("1"),
						},
					},
				},
				&elasticache.TagListMessage{
					TagList: []*elasticache.Tag{
						&elasticache.Tag{
							Key:   aws.String("test"),
							Value: aws.String("2"),
						},
					},
				},
			},
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				&serviceConnector{
					region: "test-1",
					elasticache: mockElasticCache{
						ltfro:   &elasticache.TagListMessage{},
						ltfrerr: errors.New("error with test-1"),
					},
				},
				&serviceConnector{
					region: "test-2",
					elasticache: mockElasticCache{
						ltfro: &elasticache.TagListMessage{
							TagList: []*elasticache.Tag{
								&elasticache.Tag{
									Key:   aws.String("test"),
									Value: aws.String("2"),
								}},
						},
						ltfrerr: nil,
					},
				},
			},
			expectedError: Errs{
				&callErr{
					err:     errors.New("error with test-1"),
					region:  "test-1",
					service: elasticache.ServiceName,
				},
			},
			expectedTags: []*elasticache.TagListMessage{
				&elasticache.TagListMessage{
					TagList: []*elasticache.Tag{
						&elasticache.Tag{
							Key:   aws.String("test"),
							Value: aws.String("2"),
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		c := &Connector{svcs: tt.mocked}
		tags, err := c.GetElasticacheTags(nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(tags, tt.expectedTags) {
			t.Errorf("%s [%d] - tags: received=%+v | expected=%+v",
				tt.name, i, tags, tt.expectedTags)
		}
	}
}
