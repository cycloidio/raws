package raws

import (
	"errors"
	"testing"

	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type mockS3 struct {
	s3iface.S3API

	// Mock of ListBucket
	lbo   *s3.ListBucketsOutput
	lberr error

	// Mock of DescribeVpcs
	gbto   *s3.GetBucketTaggingOutput
	gbterr error

	// Mock of DescribeImages
	loo   *s3.ListObjectsOutput
	loerr error

	// Mock of DescribeSecurityGroups
	gotout *s3.GetObjectTaggingOutput
	goterr error
}

func (m mockS3) ListBuckets(*s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	return m.lbo, m.lberr
}

func (m mockS3) GetBucketTagging(*s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error) {
	return m.gbto, m.gbterr
}

func (m mockS3) ListObjects(*s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return m.loo, m.loerr
}

func (m mockS3) GetObjectTagging(*s3.GetObjectTaggingInput) (*s3.GetObjectTaggingOutput, error) {
	return m.gotout, m.goterr
}

func TestGetBuckets(t *testing.T) {
	tests := []struct {
		name            string
		mocked          []*serviceConnector
		expectedBuckets []*s3.ListBucketsOutput
		expectedError   Errs
	}{{name: "one region with error",
		mocked: []*serviceConnector{
			{
				region: "test",
				s3: mockS3{
					lbo:   &s3.ListBucketsOutput{},
					lberr: errors.New("error with test"),
				},
			},
		},
		expectedError: Errs{&callErr{
			err:     errors.New("error with test"),
			region:  "test",
			service: s3.ServiceName,
		}},
		expectedBuckets: nil,
	},
		{name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					s3: mockS3{
						lbo: &s3.ListBucketsOutput{
							Buckets: []*s3.Bucket{{
								Name: aws.String("test"),
							}},
						},
						lberr: nil,
					},
				},
			},
			expectedError: nil,
			expectedBuckets: []*s3.ListBucketsOutput{
				{
					Buckets: []*s3.Bucket{{
						Name: aws.String("test"),
					}},
				},
			},
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					s3: mockS3{
						lbo: &s3.ListBucketsOutput{
							Buckets: []*s3.Bucket{{
								Name: aws.String("test-1"),
							}},
						},
						lberr: nil,
					},
				},
				{
					region: "test-2",
					s3: mockS3{
						lbo: &s3.ListBucketsOutput{
							Buckets: []*s3.Bucket{{
								Name: aws.String("test-2"),
							}},
						},
						lberr: nil,
					},
				},
			},
			expectedError: nil,
			expectedBuckets: []*s3.ListBucketsOutput{
				{
					Buckets: []*s3.Bucket{{
						Name: aws.String("test-1"),
					}},
				},
				{
					Buckets: []*s3.Bucket{{
						Name: aws.String("test-2"),
					}},
				},
			},
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					s3: mockS3{
						lbo:   &s3.ListBucketsOutput{},
						lberr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					s3: mockS3{
						lbo: &s3.ListBucketsOutput{
							Buckets: []*s3.Bucket{{
								Name: aws.String("test-2"),
							}},
						},
						lberr: nil,
					},
				},
			},
			expectedError: Errs{&callErr{
				err:     errors.New("error with test-1"),
				region:  "test-1",
				service: s3.ServiceName,
			}},
			expectedBuckets: []*s3.ListBucketsOutput{
				{
					Buckets: []*s3.Bucket{{
						Name: aws.String("test-2"),
					}},
				},
			},
		}}

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		buckets, err := c.ListBuckets(nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(buckets, tt.expectedBuckets) {
			t.Errorf("%s [%d] - S3 buckets: received=%+v | expected=%+v",
				tt.name, i, buckets, tt.expectedBuckets)
		}
	}
}

func TestGetBucketTags(t *testing.T) {
	tests := []struct {
		name          string
		mocked        []*serviceConnector
		expectedTags  []*s3.GetBucketTaggingOutput
		expectedError Errs
	}{{name: "one region with error",
		mocked: []*serviceConnector{
			{
				region: "test",
				s3: mockS3{
					gbto:   &s3.GetBucketTaggingOutput{},
					gbterr: errors.New("error with test"),
				},
			},
		},
		expectedError: Errs{&callErr{
			err:     errors.New("error with test"),
			region:  "test",
			service: s3.ServiceName,
		}},
		expectedTags: nil,
	},
		{name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					s3: mockS3{
						gbto: &s3.GetBucketTaggingOutput{
							TagSet: []*s3.Tag{{
								Key:   aws.String("test"),
								Value: aws.String("1"),
							}},
						},
						gbterr: nil,
					},
				},
			},
			expectedError: nil,
			expectedTags: []*s3.GetBucketTaggingOutput{
				{
					TagSet: []*s3.Tag{{
						Key:   aws.String("test"),
						Value: aws.String("1"),
					}},
				},
			},
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					s3: mockS3{
						gbto: &s3.GetBucketTaggingOutput{
							TagSet: []*s3.Tag{{
								Key:   aws.String("test"),
								Value: aws.String("1"),
							}},
						},
						gbterr: nil,
					},
				},
				{
					region: "test-2",
					s3: mockS3{
						gbto: &s3.GetBucketTaggingOutput{
							TagSet: []*s3.Tag{{
								Key:   aws.String("test"),
								Value: aws.String("2"),
							}},
						},
						gbterr: nil,
					},
				},
			},
			expectedError: nil,
			expectedTags: []*s3.GetBucketTaggingOutput{
				{
					TagSet: []*s3.Tag{{
						Key:   aws.String("test"),
						Value: aws.String("1"),
					}},
				},
				{
					TagSet: []*s3.Tag{{
						Key:   aws.String("test"),
						Value: aws.String("2"),
					}},
				},
			},
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					s3: mockS3{
						gbto:   &s3.GetBucketTaggingOutput{},
						gbterr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					s3: mockS3{
						gbto: &s3.GetBucketTaggingOutput{
							TagSet: []*s3.Tag{{
								Key:   aws.String("test"),
								Value: aws.String("2"),
							}},
						},
						gbterr: nil,
					},
				},
			},
			expectedError: Errs{&callErr{
				err:     errors.New("error with test-1"),
				region:  "test-1",
				service: s3.ServiceName,
			}},
			expectedTags: []*s3.GetBucketTaggingOutput{
				{
					TagSet: []*s3.Tag{{
						Key:   aws.String("test"),
						Value: aws.String("2"),
					}},
				},
			},
		}}

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		tags, err := c.GetBucketTags(nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(tags, tt.expectedTags) {
			t.Errorf("%s [%d] - S3 buckets tags: received=%+v | expected=%+v",
				tt.name, i, tags, tt.expectedTags)
		}
	}
}

func TestGetObjets(t *testing.T) {
	tests := []struct {
		name            string
		mocked          []*serviceConnector
		expectedObjects []*s3.ListObjectsOutput
		expectedError   Errs
	}{{name: "one region with error",
		mocked: []*serviceConnector{
			{
				region: "test",
				s3: mockS3{
					loo:   &s3.ListObjectsOutput{},
					loerr: errors.New("error with test"),
				},
			},
		},
		expectedError: Errs{&callErr{
			err:     errors.New("error with test"),
			region:  "test",
			service: s3.ServiceName,
		}},
		expectedObjects: nil,
	},
		{name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					s3: mockS3{
						loo: &s3.ListObjectsOutput{
							Name: aws.String("test"),
						},
						loerr: nil,
					},
				},
			},
			expectedError: nil,
			expectedObjects: []*s3.ListObjectsOutput{
				{
					Name: aws.String("test"),
				},
			},
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					s3: mockS3{
						loo: &s3.ListObjectsOutput{
							Name: aws.String("test-1"),
						},
						loerr: nil,
					},
				},
				{
					region: "test-2",
					s3: mockS3{
						loo: &s3.ListObjectsOutput{
							Name: aws.String("test-2"),
						},
						loerr: nil,
					},
				},
			},
			expectedError: nil,
			expectedObjects: []*s3.ListObjectsOutput{
				{
					Name: aws.String("test-1"),
				},
				{
					Name: aws.String("test-2"),
				},
			},
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					s3: mockS3{
						loo:   &s3.ListObjectsOutput{},
						loerr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					s3: mockS3{
						loo: &s3.ListObjectsOutput{
							Name: aws.String("test-2")},
						loerr: nil,
					},
				},
			},
			expectedError: Errs{&callErr{
				err:     errors.New("error with test-1"),
				region:  "test-1",
				service: s3.ServiceName,
			}},
			expectedObjects: []*s3.ListObjectsOutput{
				{
					Name: aws.String("test-2"),
				},
			},
		}}

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		objects, err := c.ListObjects(nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(objects, tt.expectedObjects) {
			t.Errorf("%s [%d] - S3 objects: received=%+v | expected=%+v",
				tt.name, i, objects, tt.expectedObjects)
		}
	}
}

func TestObjectsTags(t *testing.T) {
	tests := []struct {
		name          string
		mocked        []*serviceConnector
		expectedTags  []*s3.GetObjectTaggingOutput
		expectedError Errs
	}{{name: "one region with error",
		mocked: []*serviceConnector{
			{
				region: "test",
				s3: mockS3{
					gotout: &s3.GetObjectTaggingOutput{},
					goterr: errors.New("error with test"),
				},
			},
		},
		expectedError: Errs{&callErr{
			err:     errors.New("error with test"),
			region:  "test",
			service: s3.ServiceName,
		}},
		expectedTags: nil,
	},
		{name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					s3: mockS3{
						gotout: &s3.GetObjectTaggingOutput{
							VersionId: aws.String("test"),
						},
						goterr: nil,
					},
				},
			},
			expectedError: nil,
			expectedTags: []*s3.GetObjectTaggingOutput{
				{
					VersionId: aws.String("test"),
				},
			},
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					s3: mockS3{
						gotout: &s3.GetObjectTaggingOutput{
							VersionId: aws.String("test-1"),
						},
						goterr: nil,
					},
				},
				{
					region: "test-2",
					s3: mockS3{
						gotout: &s3.GetObjectTaggingOutput{
							VersionId: aws.String("test-2"),
						},
						goterr: nil,
					},
				},
			},
			expectedError: nil,
			expectedTags: []*s3.GetObjectTaggingOutput{
				{
					VersionId: aws.String("test-1"),
				},
				{
					VersionId: aws.String("test-2"),
				},
			},
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					s3: mockS3{
						gotout: &s3.GetObjectTaggingOutput{},
						goterr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					s3: mockS3{
						gotout: &s3.GetObjectTaggingOutput{
							VersionId: aws.String("test-2"),
						},
						goterr: nil,
					},
				},
			},
			expectedError: Errs{&callErr{
				err:     errors.New("error with test-1"),
				region:  "test-1",
				service: s3.ServiceName,
			}},
			expectedTags: []*s3.GetObjectTaggingOutput{
				{
					VersionId: aws.String("test-2"),
				},
			},
		}}

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		tags, err := c.GetObjectsTags(nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(tags, tt.expectedTags) {
			t.Errorf("%s [%d] - S3 object tags: received=%+v | expected=%+v",
				tt.name, i, tags, tt.expectedTags)
		}
	}
}
