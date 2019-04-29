package main

var (
	// functions is the list of fuctions that will be added
	// to the AWSReader with the corresponding implementation
	functions = []Function{
		// ec2
		Function{
			Entity:  "Instances",
			Prefix:  "Describe",
			Service: "ec2",
			Documentation: `
			// GetInstances returns all EC2 instances based on the input given.
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			Entity:  "Vpcs",
			Prefix:  "Describe",
			Service: "ec2",
			Documentation: `
			// GetVpcs returns all EC2 VPCs based on the input given.
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			Entity:  "Images",
			Prefix:  "Describe",
			Service: "ec2",
			Documentation: `
			// GetImages returns all EC2 AMI based on the input given.
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			Entity:        "Images",
			Prefix:        "Describe",
			Service:       "ec2",
			FilterByOwner: "Owners",
			Documentation: `
			// GetOwnImages returns all EC2 AMI belonging to the Account ID based on the input given.
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			Entity:  "SecurityGroups",
			Prefix:  "Describe",
			Service: "ec2",
			Documentation: `
			// GetSecurityGroups returns all EC2 security groups based on the input given.
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			Entity:  "Subnets",
			Prefix:  "Describe",
			Service: "ec2",
			Documentation: `
			// GetSubnets returns all EC2 subnets based on the input given.
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			Entity:  "Volumes",
			Prefix:  "Describe",
			Service: "ec2",
			Documentation: `
			// GetVolumes returns all EC2 volumes based on the input given.
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			Entity:  "Snapshots",
			Prefix:  "Describe",
			Service: "ec2",
			Documentation: `
			// GetSnapshots returns all snapshots based on the input given.
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			Entity:        "Snapshots",
			Prefix:        "Describe",
			Service:       "ec2",
			FilterByOwner: "OwnerIds",
			Documentation: `
			// GetOwnSnapshots returns all snapshots belonging to the Account ID based on the input given.
			// Returned values are commented in the interface doc comment block.
			`,
		},

		// elasticache
		Function{
			FnName:  "GetElastiCacheClusters",
			Entity:  "CacheClusters",
			Prefix:  "Describe",
			Service: "elasticache",
			Documentation: `
			// GetElastiCacheClusters returns all Elasticache clusters based on the input given.
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			FnName:   "GetElastiCacheTags",
			Entity:   "TagsForResource",
			Prefix:   "List",
			Service:  "elasticache",
			FnOutput: "TagListMessage",
			Documentation: `
			// GetElastiCacheTags returns a list of tags of Elasticache resources based on its ARN.
			// Returned values are commented in the interface doc comment block.
			`,
		},

		// elb
		Function{
			Entity:  "LoadBalancers",
			Prefix:  "Describe",
			Service: "elb",
			Documentation: `
			// GetLoadBalancers returns a list of ELB (v1) based on the input from the different regions.
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			FnName:  "GetLoadBalancersTags",
			Entity:  "Tags",
			Prefix:  "Describe",
			Service: "elb",
			Documentation: `
			// GetLoadBalancersTags returns a list of Tags based on the input from the different regions.
			// Returned values are commented in the interface doc comment block.
			`,
		},

		// elbv2
		Function{
			FnName:  "GetLoadBalancersV2",
			Entity:  "LoadBalancers",
			Prefix:  "Describe",
			Service: "elbv2",
			Documentation: `
			// GetLoadBalancersV2 returns a list of ELB (v2) - also known as ALB - based on the input from the different regions.
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			FnName:  "GetLoadBalancersV2Tags",
			Entity:  "Tags",
			Prefix:  "Describe",
			Service: "elbv2",
			Documentation: `
			// GetLoadBalancersV2Tags returns a list of Tags based on the input from the different regions.
			// Returned values are commented in the interface doc comment block.
			`,
		},

		// rds
		Function{
			Entity:  "DBInstances",
			Prefix:  "Describe",
			Service: "rds",
			Documentation: `
			// GetDBInstances returns all DB instances based on the input given.
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			FnName:  "GetDBInstancesTags",
			Entity:  "TagsForResource",
			Prefix:  "List",
			Service: "rds",
			Documentation: `
			// GetDBInstancesTags returns a list of tags from an ARN, extra filters for tags can also be provided.
			// Returned values are commented in the interface doc comment block.
			`,
		},

		// s3
		Function{
			// TODO: https://github.com/cycloidio/raws/issues/44
			FnName:  "ListBuckets",
			Entity:  "Buckets",
			Prefix:  "List",
			Service: "s3",
			Documentation: `
			// ListBuckets returns all S3 buckets based on the input given.
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			FnName:  "GetBucketTags",
			Entity:  "BucketTagging",
			Prefix:  "Get",
			Service: "s3",
			Documentation: `
			// GetBucketTags returns tags associated with S3 buckets based on the input given.
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			// TODO: https://github.com/cycloidio/raws/issues/44
			FnName:  "ListObjects",
			Entity:  "Objects",
			Prefix:  "List",
			Service: "s3",
			Documentation: `
			// ListObjects returns a list of all S3 objects in a bucket based on the input given.
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			FnName:  "GetObjectsTags",
			Entity:  "ObjectTagging",
			Prefix:  "Get",
			Service: "s3",
			Documentation: `
			// GetObjectsTags returns tags associated with S3 objects based on the input given.
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			FnName:  "GetRecordedResourceCounts",
			Entity:  "DiscoveredResourceCounts",
			Prefix:  "Get",
			Service: "configservice",
			Documentation: `
			// GetRecordedResourceCounts returns counts of the AWS resources which have
			// been recorded by AWS Config.
			// See https://docs.aws.amazon.com/config/latest/APIReference/API_GetDiscoveredResourceCounts.html
			// for more information about what to enable in your AWS account, the list of
			// supported resources, etc.
			`,
		},

		// s3downloader
		Function{
			FnSignature:  "DownloadObject(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error)",
			NoGenerateFn: true,
			Documentation: `
			// DownloadObject downloads an object in a bucket based on the input given
			`,
		},

		// cloudfront
		Function{
			FnName:  "GetCloudFrontDistributions",
			Entity:  "Distributions",
			Prefix:  "List",
			Service: "cloudfront",
			Documentation: `
			// GetCloudFrontDistributions returns all the CloudFront Distributions on the given input
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			FnName:  "GetCloudFrontPublicKeys",
			Entity:  "PublicKeys",
			Prefix:  "List",
			Service: "cloudfront",
			Documentation: `
			// GetCloudFrontPublicKeys returns all the CloudFront Public Keys on the given input
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			Entity:  "CloudFrontOriginAccessIdentities",
			Prefix:  "List",
			Service: "cloudfront",
			Documentation: `
			// GetCloudFrontOriginAccessIdentities returns all the CloudFront Origin Access Identities on the given input
			// Returned values are commented in the interface doc comment block.
			`,
		},
	}
)
