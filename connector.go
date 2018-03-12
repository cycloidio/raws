package raws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/elasticache/elasticacheiface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

// AWSReader is the interface defining all methods that need to be implemented
//
// The next behavior commented in the below paragraph, applies to every method
// which clearly match what's explained, for the sake of not repeating the same,
// over and over.
// The most of the methods defined by this interface, return their results in a
// map. Those maps, have as keys, the AWS region which have been requested and
// the values are the items returned by AWS for such region.
// Because the methods may make calls to different regions, in case that there
// is an error on a region, the returned map won't have any entry for such
// region and such errors will be reported by the returned error, nonetheless
// the items, got from the successful requests to other regions, will be
// returned, with the meaning that the methods will return partial results, in
// case of errors.
// For avoiding by the callers the problem of if the returned map may be nil,
// the function will always return a map instance, which will be of length 0
// in case that there is not any successful request.
type AWSReader interface {
	// GetAccountID returns the current ID for the account used
	GetAccountID() string

	// GetRegions returns the currently used regions for the Connector
	GetRegions() []string

	// GetInstances returns all EC2 instances based on the input given.
	// Returned values are commented in the interface doc comment block.
	GetInstances(ctx context.Context, input *ec2.DescribeInstancesInput) (map[string]ec2.DescribeInstancesOutput, error)

	// GetVpcs returns all EC2 VPCs based on the input given.
	// Returned values are commented in the interface doc comment block.
	GetVpcs(ctx context.Context, input *ec2.DescribeVpcsInput) (map[string]ec2.DescribeVpcsOutput, error)

	// GetImages returns all EC2 AMI belonging to the Account ID based on the input given.
	// Returned values are commented in the interface doc comment block.
	GetImages(ctx context.Context, input *ec2.DescribeImagesInput) (map[string]ec2.DescribeImagesOutput, error)

	// GetSecurityGroups returns all EC2 security groups based on the input given.
	// Returned values are commented in the interface doc comment block.
	GetSecurityGroups(
		ctx context.Context, input *ec2.DescribeSecurityGroupsInput,
	) (map[string]ec2.DescribeSecurityGroupsOutput, error)

	// GetSubnets returns all EC2 subnets based on the input given.
	// Returned values are commented in the interface doc comment block.
	GetSubnets(ctx context.Context, input *ec2.DescribeSubnetsInput) (map[string]ec2.DescribeSubnetsOutput, error)

	// GetVolumes returns all EC2 volumes based on the input given.
	// Returned values are commented in the interface doc comment block.
	GetVolumes(ctx context.Context, input *ec2.DescribeVolumesInput) (map[string]ec2.DescribeVolumesOutput, error)

	// GetSnapshots returns all snapshots belonging to the Account ID based on the input given.
	// Returned values are commented in the interface doc comment block.
	GetSnapshots(ctx context.Context, input *ec2.DescribeSnapshotsInput) (map[string]ec2.DescribeSnapshotsOutput, error)

	// GetElastiCacheCluster returns all Elasticache clusters based on the input given.
	// Returned values are commented in the interface doc comment block.
	GetElastiCacheCluster(
		ctx context.Context, input *elasticache.DescribeCacheClustersInput,
	) (map[string]elasticache.DescribeCacheClustersOutput, error)

	// GetElastiCacheTags returns a list of tags of Elasticache resources based on its ARN.
	// Returned values are commented in the interface doc comment block.
	GetElastiCacheTags(
		ctx context.Context, input *elasticache.ListTagsForResourceInput,
	) (map[string]elasticache.TagListMessage, error)

	// GetLoadBalancers returns a list of ELB (v1) based on the input from the different regions.
	// Returned values are commented in the interface doc comment block.
	GetLoadBalancers(
		ctx context.Context, input *elb.DescribeLoadBalancersInput,
	) (map[string]elb.DescribeLoadBalancersOutput, error)

	// GetLoadBalancersTags returns a list of Tags based on the input from the different regions.
	// Returned values are commented in the interface doc comment block.
	GetLoadBalancersTags(ctx context.Context, input *elb.DescribeTagsInput) (map[string]elb.DescribeTagsOutput, error)

	// GetLoadBalancersV2 returns a list of ELB (v2) - also known as ALB - based on the input from the different regions.
	// Returned values are commented in the interface doc comment block.
	GetLoadBalancersV2(
		ctx context.Context, input *elbv2.DescribeLoadBalancersInput,
	) (map[string]elbv2.DescribeLoadBalancersOutput, error)

	// GetLoadBalancersV2Tags returns a list of Tags based on the input from the different regions.
	// Returned values are commented in the interface doc comment block.
	GetLoadBalancersV2Tags(
		ctx context.Context, input *elbv2.DescribeTagsInput,
	) (map[string]elbv2.DescribeTagsOutput, error)

	// GetDBInstances returns all DB instances based on the input given.
	// Returned values are commented in the interface doc comment block.
	GetDBInstances(
		ctx context.Context, input *rds.DescribeDBInstancesInput,
	) (map[string]rds.DescribeDBInstancesOutput, error)

	// GetDBInstancesTags returns a list of tags from an ARN, extra filters for tags can also be provided.
	// Returned values are commented in the interface doc comment block.
	GetDBInstancesTags(
		ctx context.Context, input *rds.ListTagsForResourceInput,
	) (map[string]rds.ListTagsForResourceOutput, error)

	// ListBuckets returns all S3 buckets based on the input given.
	// Returned values are commented in the interface doc comment block.
	ListBuckets(ctx context.Context, input *s3.ListBucketsInput) (map[string]s3.ListBucketsOutput, error)

	// GetBucketTags returns tags associated with S3 buckets based on the input given.
	// Returned values are commented in the interface doc comment block.
	GetBucketTags(ctx context.Context, input *s3.GetBucketTaggingInput) (map[string]s3.GetBucketTaggingOutput, error)

	// ListObjects returns a list of all S3 objects in a bucket based on the input given.
	// Returned values are commented in the interface doc comment block.
	ListObjects(ctx context.Context, input *s3.ListObjectsInput) (map[string]s3.ListObjectsOutput, error)

	// DownloadObject downloads an object in a bucket based on the input given
	DownloadObject(
		ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader),
	) (int64, error)

	// GetObjectsTags returns tags associated with S3 objects based on the input given.
	// Returned values are commented in the interface doc comment block.
	GetObjectsTags(ctx context.Context, input *s3.GetObjectTaggingInput) (map[string]s3.GetObjectTaggingOutput, error)

	// GetRecordedResourceCounts returns counts of the AWS resources which have
	// been recorded by AWS Config.
	// See https://docs.aws.amazon.com/config/latest/APIReference/API_GetDiscoveredResourceCounts.html
	// for more information about what to enable in your AWS account, the list of
	// supported resources, etc.
	GetRecordedResourceCounts(
		ctx context.Context, input *configservice.GetDiscoveredResourceCountsInput,
	) (map[string]configservice.GetDiscoveredResourceCountsOutput, error)
}

// NewAWSReader returns an object which also contains the accountID and extend the different regions to use.
//
// The accountID is helpful to return only the AMI or snapshots that belong to the account.
//
// While the regions slice also supports regex so, "eu-*" can be passed, and will be extended to: eu-west-1, eu-west-2 &
// eu-central-1.
//
// When calls are done through the Connector instance, then all regions will be called for those services.
// Thus making requests to AWS much easier than through the different connectors/regions of its go SDK.
//
// The connections are not all established while instancing, but the various sessions are, this way connections are only
// made for services that are called, otherwise only the sessions remain.
func NewAWSReader(
	ctx context.Context, accessKey string, secretKey string, regions []string, config *aws.Config,
) (AWSReader, error) {
	var c = connector{}

	creds, ec2s, sts, err := configureAWS(accessKey, secretKey)
	if err != nil {
		return nil, err
	}
	c.creds = creds
	if err := c.setAccountID(ctx, sts); err != nil {
		return nil, err
	}
	if err := c.setRegions(ctx, ec2s, regions); err != nil {
		return nil, err
	}
	c.setServices(config)
	return &c, nil
}

// The connector provides easy access to AWS SDK calls.
//
// By using it, calls can be made directly through multiple regions, and will filter only data that belongs to you.
// For example, when fetching the list of AMI, or snapshots.
//
// In order to start making calls, only calling NewAWSReader is required.
type connector struct {
	regions   []string
	svcs      []*serviceConnector
	creds     *credentials.Credentials
	accountID *string
}

func (c *connector) GetAccountID() string {
	return *c.accountID
}

func (c *connector) GetRegions() []string {
	return c.regions
}

type serviceConnector struct {
	region       string
	session      *session.Session
	ec2          ec2iface.EC2API
	elb          elbiface.ELBAPI
	elbv2        elbv2iface.ELBV2API
	rds          rdsiface.RDSAPI
	s3           s3iface.S3API
	s3downloader s3manageriface.DownloaderAPI
	elasticache  elasticacheiface.ElastiCacheAPI
}

func configureAWS(accessKey string, secretKey string) (*credentials.Credentials, ec2iface.EC2API, stsiface.STSAPI, error) {
	/* The default region is only used to (1) get the list of region and
	 * (2) get the account ID associated with the credentials.
	 *
	 * It is not used as a default region for services, therefore if no
	 * region is specified when instantiating the connector, then it will
	 * not try to establish any connections with AWS services.
	 */
	const defaultRegion string = "eu-west-1"
	var token = ""

	creds := credentials.NewStaticCredentials(accessKey, secretKey, token)
	_, err := creds.Get()
	if err != nil {
		return nil, nil, nil, err
	}
	sess := session.Must(
		session.NewSession(&aws.Config{
			Region:      aws.String(defaultRegion),
			DisableSSL:  aws.Bool(false),
			MaxRetries:  aws.Int(3),
			Credentials: creds,
		}),
	)
	return creds, ec2.New(sess), sts.New(sess), nil
}

func (c *connector) setRegions(ctx context.Context, ec2 ec2iface.EC2API, enabledRegions []string) error {
	if len(enabledRegions) == 0 {
		return errors.New("at least one region name is required")
	}
	regions, err := ec2.DescribeRegionsWithContext(ctx, nil)
	if err != nil {
		return err
	}
	for _, enabledRegion := range enabledRegions {
		for _, region := range regions.Regions {
			if match, _ := filepath.Match(enabledRegion, *region.RegionName); match {
				c.regions = append(c.regions, *region.RegionName)
			}
		}
	}
	if len(c.regions) == 0 {
		return fmt.Errorf("found 0 regions matching: %v", enabledRegions)
	}
	return nil
}

func (c *connector) setAccountID(ctx context.Context, sts stsiface.STSAPI) error {
	resp, err := sts.GetCallerIdentityWithContext(ctx, nil)
	if err != nil {
		return err
	}
	c.accountID = resp.Account
	return nil
}

func (c *connector) setServices(config *aws.Config) {
	if config != nil {
		config.Credentials = c.creds
	} else {
		config = &aws.Config{
			DisableSSL:  aws.Bool(false),
			MaxRetries:  aws.Int(3),
			Credentials: c.creds,
		}
	}
	for _, region := range c.regions {
		config.Region = aws.String(region)
		sess := session.Must(session.NewSession(config))
		svc := &serviceConnector{
			region:  region,
			session: sess,
		}
		c.svcs = append(c.svcs, svc)
	}
}
