package raws

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Returns all S3 buckets based on the input given
func (c *Connector) ListBuckets(input *s3.ListBucketsInput) ([]*s3.ListBucketsOutput, Errs) {
	var errs Errs
	var bucketsList []*s3.ListBucketsOutput

	for _, svc := range c.svcs {
		if svc.s3 == nil {
			svc.s3 = s3.New(svc.session)
		}
		buckets, err := svc.s3.ListBuckets(input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, s3.ServiceName, err))
		} else {
			bucketsList = append(bucketsList, buckets)
		}
	}
	return bucketsList, errs
}

// Returns tags associated with S3 buckets based on the input given
func (c *Connector) GetBucketTags(input *s3.GetBucketTaggingInput) ([]*s3.GetBucketTaggingOutput, Errs) {
	var errs Errs
	var bucketsTagList []*s3.GetBucketTaggingOutput

	for _, svc := range c.svcs {
		if svc.s3 == nil {
			svc.s3 = s3.New(svc.session)
		}
		bucketsTags, err := svc.s3.GetBucketTagging(input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, s3.ServiceName, err))
		} else {
			bucketsTagList = append(bucketsTagList, bucketsTags)
		}
	}
	return bucketsTagList, errs
}

// Returns a list of all S3 objects in a bucket based on the input given
func (c *Connector) ListObjects(input *s3.ListObjectsInput) ([]*s3.ListObjectsOutput, Errs) {
	var errs Errs
	var objectsList []*s3.ListObjectsOutput

	for _, svc := range c.svcs {
		if svc.s3 == nil {
			svc.s3 = s3.New(svc.session)
		}
		buckets, err := svc.s3.ListObjects(input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, s3.ServiceName, err))
		} else {
			objectsList = append(objectsList, buckets)
		}
	}
	return objectsList, errs
}

// DownloadObject downloads an object in a bucket based on the input given
func (c *Connector) DownloadObject(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error) {
	var err error = nil
	var n int64 = 0

	for _, svc := range c.svcs {
		if svc.s3downloader == nil {
			svc.s3downloader = s3manager.NewDownloader(svc.session)
		}
		n, err = svc.s3downloader.Download(w, input, options...)
		if err == nil {
			return n, nil
		}
	}
	return n, fmt.Errorf("Couldn't download '%v' in any of '%v' regions", input, c.GetRegions())
}

// Returns tags associated with S3 objects based on the input given
func (c *Connector) GetObjectsTags(input *s3.GetObjectTaggingInput) ([]*s3.GetObjectTaggingOutput, Errs) {
	var errs Errs
	var objectsTagsList []*s3.GetObjectTaggingOutput

	for _, svc := range c.svcs {
		if svc.s3 == nil {
			svc.s3 = s3.New(svc.session)
		}
		buckets, err := svc.s3.GetObjectTagging(input)
		if err != nil {
			errs = append(errs, NewAPIError(svc.region, s3.ServiceName, err))
		} else {
			objectsTagsList = append(objectsTagsList, buckets)
		}
	}
	return objectsTagsList, errs
}
