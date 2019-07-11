package raws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (c *connector) ListBuckets(ctx context.Context, input *s3.ListBucketsInput) (map[string]s3.ListBucketsOutput, error) {
	var errs Errors
	var regionsOpts = map[string]s3.ListBucketsOutput{}

	for _, svc := range c.svcs {
		if svc.s3 == nil {
			svc.s3 = s3.New(svc.session)
		}

		opt, err := svc.s3.ListBucketsWithContext(ctx, input)
		if err != nil {
			errs = append(errs, NewError(svc.region, s3.ServiceName, err))
		} else {
			newOpt := s3.ListBucketsOutput{
				Owner:   opt.Owner,
				Buckets: make([]*s3.Bucket, 0),
			}
			for _, bucket := range opt.Buckets {
				inputLocation := &s3.GetBucketLocationInput{
					Bucket: bucket.Name,
				}
				result, err := svc.s3.GetBucketLocation(inputLocation)
				if err != nil {
					errs = append(errs, NewError(svc.region, s3.ServiceName, err))
				}
				if s3.NormalizeBucketLocation(aws.StringValue(result.LocationConstraint)) == svc.region {
					newOpt.Buckets = append(newOpt.Buckets, bucket)
				}
			}
			regionsOpts[svc.region] = newOpt
		}
	}

	if errs != nil {
		return regionsOpts, errs
	}

	return regionsOpts, nil
}
