package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/elasticache/elasticacheiface"
	"testing"
	"reflect"
	"errors"
)

type mockElasticCache struct {
	elasticacheiface.ElastiCacheAPI

	// Mocking of DescribeCacheClusters
	dcco  *elasticache.DescribeCacheClustersOutput
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

func CheckErrors(t *testing.T, err error, expected error) {
	if err != nil && !reflect.DeepEqual(err, expected) {
		t.Errorf("Error received: '%v' expected '%v'",
			err.Error(), expected.Error())
	}
}
func TestGetElasticacheCluster(t *testing.T) {
	tests := []struct {
		name          string
		mocked        []*serviceConnector
		expectedClusters  []*elasticache.DescribeCacheClustersOutput
		expectedError error
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
						dcco: &elasticache.DescribeCacheClustersOutput{},
						dccerr: errors.New("error with test"),
					},
				},
			},
			expectedError: RawsErr{
				APIErrs: []callErr{
					{
						err: errors.New("error with test"),
						region: "test",
						service: elasticache.ServiceName,
					},
				},
			},
			expectedClusters: []*elasticache.DescribeCacheClustersOutput{
				&elasticache.DescribeCacheClustersOutput{},
			},
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
						dcco: &elasticache.DescribeCacheClustersOutput{},
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
			expectedError: RawsErr{
				APIErrs: []callErr{
					{
						err: errors.New("error with test"),
						region: "test-1",
						service: elasticache.ServiceName,
					},
				},
			},
			expectedClusters: []*elasticache.DescribeCacheClustersOutput{
				&elasticache.DescribeCacheClustersOutput{},
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

	for _, tt := range tests {
		c := &Connector{svcs: tt.mocked}
		tags, err := c.GetElasticCacheCluster(nil)
		CheckErrors(t, err, tt.expectedError)
		if !reflect.DeepEqual(tags, tt.expectedClusters) {
			t.Errorf("Clusters received: '%v' expected '%v'",
				tags, tt.expectedClusters)
		}
	}
}

func TestGetElasticacheTags(t *testing.T) {
	tests := []struct {
		name          string
		mocked        []*serviceConnector
		expectedTags  []*elasticache.TagListMessage
		expectedError error
	}{{name: "one region no error",
		mocked: []*serviceConnector{
			&serviceConnector{
				region: "test",
				elasticache: mockElasticCache{
					ltfro: &elasticache.TagListMessage{
						TagList: []*elasticache.Tag{
							&elasticache.Tag{
								Key: aws.String("test"),
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
						Key: aws.String("test"),
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
			expectedError: RawsErr{
				APIErrs: []callErr{
					{
						err: errors.New("error with test"),
						region: "test",
						service: elasticache.ServiceName,
					},
				},
			},
			expectedTags: []*elasticache.TagListMessage{
				&elasticache.TagListMessage{
					TagList: []*elasticache.Tag{},
				},
			},
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				&serviceConnector{
					region: "test",
					elasticache: mockElasticCache{
						ltfro: &elasticache.TagListMessage{
							TagList: []*elasticache.Tag{
								&elasticache.Tag{
									Key: aws.String("test"),
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
									Key: aws.String("test"),
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
							Key: aws.String("test"),
							Value: aws.String("1"),
						},
					},
				},
				&elasticache.TagListMessage{
					TagList: []*elasticache.Tag{
						&elasticache.Tag{
							Key: aws.String("test"),
							Value: aws.String("2"),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		c := &Connector{svcs: tt.mocked}
		tags, err := c.GetElasticacheTags(nil)
		CheckErrors(t, err, tt.expectedError)
		if !reflect.DeepEqual(tags, tt.expectedTags) {
			t.Errorf("Tags received: '%v' expected '%v'",
				tags, tt.expectedTags)
		}
	}
}