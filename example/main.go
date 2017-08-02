package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cycloidio/raws"
)

func callEC2(c *raws.Connector) {
	instances, _ := c.GetInstances(nil)
	fmt.Println(instances)
	vpcs, _ := c.GetVpcs(nil)
	fmt.Println(vpcs)
	sg, _ := c.GetSecurityGroups(nil)
	fmt.Println(sg)
	subnets, _ := c.GetSubnets(nil)
	fmt.Println(subnets)
	images, _ := c.GetImages(nil)
	fmt.Println(images)
	volumes, _ := c.GetVolumes(nil)
	fmt.Println(volumes)
	snapshots, _ := c.GetSnapshots(nil)
	fmt.Println(snapshots)
}

func callELB(c *raws.Connector) {
	elbs, _ := c.GetLoadBalancers(nil)
	fmt.Println(elbs)
	tags, _ := c.GetLoadBalancersTags(nil)
	fmt.Println(tags)
}

func callELBv2(c *raws.Connector) {
	elbs, _ := c.GetLoadBalancersV2(nil)
	fmt.Println(elbs)
	tags, _ := c.GetLoadBalancersV2Tags(nil)
	fmt.Println(tags)
}

func callRDS(c *raws.Connector) {
	instances, _ := c.GetDBInstances(nil)
	fmt.Println(instances)
	i := &rds.ListTagsForResourceInput{
		ResourceName: aws.String("MY_RDS_ARN"),
	}
	tags, _ := c.GetDBInstancesTags(i)
	fmt.Println(tags)
}

func callElastiCache(c *raws.Connector) {
	clusters, _ := c.GetElasticCacheCluster(nil)

	fmt.Println(clusters)
	i := &elasticache.ListTagsForResourceInput{
		ResourceName: aws.String("MY_ELASTICACHE_ARN"),
	}
	tags, _ := c.GetElasticacheTags(i)
	fmt.Println(tags)
}

func callS3(c *raws.Connector) {
	buckets, _ := c.ListBuckets(nil)
	fmt.Println(buckets)
	i := &s3.GetBucketTaggingInput{Bucket: aws.String("MY_BUCKET")}
	bucketsTags, _ := c.GetBucketTags(i)
	fmt.Println(bucketsTags)
	i2 := &s3.ListObjectsInput{Bucket: aws.String("MY_BUCKET")}
	objects, _ := c.ListObjects(i2)
	fmt.Println(objects)
	i3 := &s3.GetObjectTaggingInput{
		Bucket: aws.String("MY_BUCKET"),
		Key:    aws.String("MY_KEY"),
	}
	objectsTags, _ := c.GetObjectsTags(i3)
	fmt.Println(objectsTags)
}

func main() {
	accessKey := ""
	secretKey := ""
	region := []string{"eu-*"}

	c, err := raws.NewConnector(accessKey, secretKey, region, nil)
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
}
