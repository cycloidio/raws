package raws

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (c *connector) ListBuckets(
	ctx context.Context, input *s3.ListBucketsInput,
) ([]*s3.ListBucketsOutput, error) {
	var errs Errs
	var bucketsList []*s3.ListBucketsOutput

	for _, svc := range c.svcs {
		if svc.s3 == nil {
			svc.s3 = s3.New(svc.session)
		}
		buckets, err := svc.s3.ListBucketsWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, s3.ServiceName, err))
		} else {
			bucketsList = append(bucketsList, buckets)
		}
	}

	if errs != nil {
		return bucketsList, errs
	}

	return bucketsList, nil
}

func (c *connector) GetBucketTags(
	ctx context.Context, input *s3.GetBucketTaggingInput,
) ([]*s3.GetBucketTaggingOutput, error) {
	var errs Errs
	var bucketsTagList []*s3.GetBucketTaggingOutput

	for _, svc := range c.svcs {
		if svc.s3 == nil {
			svc.s3 = s3.New(svc.session)
		}
		bucketsTags, err := svc.s3.GetBucketTaggingWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, s3.ServiceName, err))
		} else {
			bucketsTagList = append(bucketsTagList, bucketsTags)
		}
	}

	if errs != nil {
		return bucketsTagList, errs
	}

	return bucketsTagList, nil
}

func (c *connector) ListObjects(
	ctx context.Context, input *s3.ListObjectsInput,
) ([]*s3.ListObjectsOutput, Errs) {
	var errs Errs
	var objectsList []*s3.ListObjectsOutput

	for _, svc := range c.svcs {
		if svc.s3 == nil {
			svc.s3 = s3.New(svc.session)
		}
		buckets, err := svc.s3.ListObjectsWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, s3.ServiceName, err))
		} else {
			objectsList = append(objectsList, buckets)
		}
	}
	return objectsList, errs
}

func (c *connector) DownloadObject(
	ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader),
) (int64, error) {
	var err error
	var n int64

	if input.Bucket == nil || input.Key == nil {
		return n, fmt.Errorf("couldn't download undefined object (keys or bucket not set)")
	}
	for _, svc := range c.svcs {
		if svc.s3downloader == nil {
			svc.s3downloader = s3manager.NewDownloader(svc.session)
		}
		n, err = svc.s3downloader.DownloadWithContext(ctx, w, input, options...)
		if err == nil {
			return n, nil
		}
	}
	return n, fmt.Errorf("couldn't download '%s/%s' in any of '%+v' regions", *input.Bucket, *input.Key, c.GetRegions())
}

func (c *connector) GetObjectsTags(
	ctx context.Context, input *s3.GetObjectTaggingInput,
) ([]*s3.GetObjectTaggingOutput, Errs) {
	var errs Errs
	var objectsTagsList []*s3.GetObjectTaggingOutput

	for _, svc := range c.svcs {
		if svc.s3 == nil {
			svc.s3 = s3.New(svc.session)
		}
		buckets, err := svc.s3.GetObjectTaggingWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, s3.ServiceName, err))
		} else {
			objectsTagsList = append(objectsTagsList, buckets)
		}
	}
	return objectsTagsList, errs
}
