package raws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type mockS3 struct {
	s3iface.S3API

	// Mock of ListBucket
	lbo   *s3.ListBucketsOutput
	lberr error

	// Mock of GetBucketTags
	gbto   *s3.GetBucketTaggingOutput
	gbterr error

	// Mock of ListObjects
	loo   *s3.ListObjectsOutput
	loerr error

	// Mock of GetObjectTags
	gotout *s3.GetObjectTaggingOutput
	goterr error

	// Mock of DownloadObject
	dob   int64
	doerr error
}

func (m mockS3) ListBucketsWithContext(
	_ aws.Context, _ *s3.ListBucketsInput, _ ...request.Option,
) (*s3.ListBucketsOutput, error) {
	return m.lbo, m.lberr
}

func (m mockS3) GetBucketTaggingWithContext(
	_ aws.Context, _ *s3.GetBucketTaggingInput, _ ...request.Option,
) (*s3.GetBucketTaggingOutput, error) {
	return m.gbto, m.gbterr
}

type mockS3Downloader struct {
	s3manageriface.DownloaderAPI

	// Mock of DownloadObject
	dob   int64
	doerr error
}

func (m mockS3Downloader) DownloadWithContext(
	_ aws.Context, _ io.WriterAt, _ *s3.GetObjectInput, _ ...func(*s3manager.Downloader),
) (int64, error) {
	return m.dob, m.doerr
}

func (m mockS3) ListObjectsWithContext(
	_ aws.Context, _ *s3.ListObjectsInput, _ ...request.Option,
) (*s3.ListObjectsOutput, error) {
	return m.loo, m.loerr
}

func (m mockS3) GetObjectTaggingWithContext(
	_ aws.Context, _ *s3.GetObjectTaggingInput, _ ...request.Option,
) (*s3.GetObjectTaggingOutput, error) {
	return m.gotout, m.goterr
}

func TestListBuckets(t *testing.T) {
	tests := []struct {
		name            string
		mocked          []*serviceConnector
		expectedBuckets []*s3.ListBucketsOutput
		expectedError   error
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
		expectedError: Errors{Error{
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
			expectedError: Errors{Error{
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

	var ctx = context.Background()

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		buckets, err := c.ListBuckets(ctx, nil)
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
		expectedError error
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
		expectedError: Errors{Error{
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
			expectedError: Errors{Error{
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

	var ctx = context.Background()

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		tags, err := c.GetBucketTags(ctx, nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(tags, tt.expectedTags) {
			t.Errorf("%s [%d] - S3 buckets tags: received=%+v | expected=%+v",
				tt.name, i, tags, tt.expectedTags)
		}
	}
}

func TestListObjects(t *testing.T) {
	tests := []struct {
		name            string
		mocked          []*serviceConnector
		expectedObjects []*s3.ListObjectsOutput
		expectedError   error
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
		expectedError: Errors{Error{
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
			expectedError: Errors{Error{
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

	var ctx = context.Background()

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		objects, err := c.ListObjects(ctx, nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(objects, tt.expectedObjects) {
			t.Errorf("%s [%d] - S3 objects: received=%+v | expected=%+v",
				tt.name, i, objects, tt.expectedObjects)
		}
	}
}

func TestDownloadObjet(t *testing.T) {

	tests := []struct {
		name          string
		mocked        []*serviceConnector
		regions       []string
		input         *s3.GetObjectInput
		expectedBytes int64
		expectedError error
	}{
		{
			name: "one region with error",
			mocked: []*serviceConnector{
				{
					s3downloader: mockS3Downloader{
						dob:   0,
						doerr: errors.New("error with test"),
					},
				},
			},
			regions: []string{"test"},
			input: &s3.GetObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
			},
			expectedError: fmt.Errorf("couldn't download 'bucket/key' in any of '[test]' regions"),
			expectedBytes: 0,
		},
		{
			name:    "invalid file requested",
			mocked:  nil,
			regions: []string{"test"},
			input: &s3.GetObjectInput{
				Bucket: nil,
				Key:    nil,
			},
			expectedError: fmt.Errorf("couldn't download undefined object (keys or bucket not set)"),
			expectedBytes: 0,
		},
		{
			name: "one region no error",
			mocked: []*serviceConnector{
				{
					s3downloader: mockS3Downloader{
						dob:   42,
						doerr: nil,
					},
				},
			},
			regions: []string{"test"},
			input: &s3.GetObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
			},
			expectedError: nil,
			expectedBytes: 42,
		},
	}

	var ctx = context.Background()

	for i, tt := range tests {
		c := &connector{
			regions: tt.regions,
			svcs:    tt.mocked,
		}
		bytes, err := c.DownloadObject(ctx, nil, tt.input, nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if tt.expectedBytes != bytes {
			t.Errorf("%s [%d] - S3 download object: received=%+v | expected=%+v",
				tt.name, i, bytes, tt.expectedBytes)
		}
	}
}

func TestGetObjectsTags(t *testing.T) {
	tests := []struct {
		name          string
		mocked        []*serviceConnector
		expectedTags  []*s3.GetObjectTaggingOutput
		expectedError error
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
		expectedError: Errors{Error{
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
			expectedError: Errors{Error{
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

	var ctx = context.Background()

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		tags, err := c.GetObjectsTags(ctx, nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(tags, tt.expectedTags) {
			t.Errorf("%s [%d] - S3 object tags: received=%+v | expected=%+v",
				tt.name, i, tags, tt.expectedTags)
		}
	}
}
