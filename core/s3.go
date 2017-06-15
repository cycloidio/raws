package core

import (
	"github.com/aws/aws-sdk-go/service/s3"
)

// Returns all S3 buckets based on the input given
func (c *Connector) GetBuckets(input *s3.ListBucketsInput) ([]*s3.ListBucketsOutput, Errs) {
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
func (c *Connector) GetObjects(input *s3.ListObjectsInput) ([]*s3.ListObjectsOutput, Errs) {
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
