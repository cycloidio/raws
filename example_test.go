package raws_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cycloidio/raws"
)

func callEC2(c raws.AWSReader) {
	var ctx = context.Background()

	instances, _ := c.GetInstances(ctx, nil)
	fmt.Println(instances)

	vpcs, _ := c.GetVpcs(ctx, nil)
	fmt.Println(vpcs)

	sg, _ := c.GetSecurityGroups(ctx, nil)
	fmt.Println(sg)

	subnets, _ := c.GetSubnets(ctx, nil)
	fmt.Println(subnets)

	images, _ := c.GetImages(ctx, nil)
	fmt.Println(images)

	volumes, _ := c.GetVolumes(ctx, nil)
	fmt.Println(volumes)

	snapshots, _ := c.GetSnapshots(ctx, nil)
	fmt.Println(snapshots)
}

func callELB(c raws.AWSReader) {
	var ctx = context.Background()

	elbs, _ := c.GetLoadBalancers(ctx, nil)
	fmt.Println(elbs)

	tags, _ := c.GetLoadBalancersTags(ctx, nil)
	fmt.Println(tags)
}

func callELBv2(c raws.AWSReader) {
	var ctx = context.Background()

	elbs, _ := c.GetLoadBalancersV2(ctx, nil)
	fmt.Println(elbs)

	tags, _ := c.GetLoadBalancersV2Tags(ctx, nil)
	fmt.Println(tags)
}

func callRDS(c raws.AWSReader) {
	var ctx = context.Background()

	instances, _ := c.GetDBInstances(ctx, nil)
	fmt.Println(instances)

	i := &rds.ListTagsForResourceInput{
		ResourceName: aws.String("MY_RDS_ARN"),
	}
	tags, _ := c.GetDBInstancesTags(ctx, i)
	fmt.Println(tags)
}

func callElastiCache(c raws.AWSReader) {
	var ctx = context.Background()

	clusters, _ := c.GetElastiCacheClusters(ctx, nil)
	fmt.Println(clusters)

	i := &elasticache.ListTagsForResourceInput{
		ResourceName: aws.String("MY_ELASTICACHE_ARN"),
	}
	tags, _ := c.GetElastiCacheTags(ctx, i)
	fmt.Println(tags)
}

func callS3(c raws.AWSReader) {
	var ctx = context.Background()

	buckets, _ := c.ListBuckets(ctx, nil)
	fmt.Println(buckets)

	i := &s3.GetBucketTaggingInput{Bucket: aws.String("MY_BUCKET")}
	bucketsTags, _ := c.GetBucketTags(ctx, i)
	fmt.Println(bucketsTags)

	i2 := &s3.ListObjectsInput{Bucket: aws.String("MY_BUCKET")}
	objects, _ := c.ListObjects(ctx, i2)
	fmt.Println(objects)

	i3 := &s3.GetObjectTaggingInput{
		Bucket: aws.String("MY_BUCKET"),
		Key:    aws.String("MY_KEY"),
	}
	objectsTags, _ := c.GetObjectsTags(ctx, i3)
	fmt.Println(objectsTags)
}

func Example() {
	var accessKey, secretKey, err = getAWSKeys()
	if err != nil {
		// When no keys the example doesn't run, but if the env vars are set
		// then it runs and make the AWS SDK calls.
		return
	}

	region := []string{"eu-*"}

	c, err := raws.NewAWSReader(context.Background(), accessKey, secretKey, region, nil)
	if err != nil {
		fmt.Printf("Error while getting NewConnector: %s\n", err.Error())
		return
	}

	callEC2(c)
	callELB(c)
	callELBv2(c)
	callRDS(c)
	callElastiCache(c)
	callS3(c)

	// Output
}
