package raws

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
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
type AWSReader interface {
	GetAccountID() string
	GetRegions() []string
	GetInstances(input *ec2.DescribeInstancesInput) ([]*ec2.DescribeInstancesOutput, Errs)
	GetVpcs(input *ec2.DescribeVpcsInput) ([]*ec2.DescribeVpcsOutput, Errs)
	GetImages(input *ec2.DescribeImagesInput) ([]*ec2.DescribeImagesOutput, Errs)
	GetSecurityGroups(input *ec2.DescribeSecurityGroupsInput) ([]*ec2.DescribeSecurityGroupsOutput, Errs)
	GetSubnets(input *ec2.DescribeSubnetsInput) ([]*ec2.DescribeSubnetsOutput, Errs)
	GetVolumes(input *ec2.DescribeVolumesInput) ([]*ec2.DescribeVolumesOutput, Errs)
	GetSnapshots(input *ec2.DescribeSnapshotsInput) ([]*ec2.DescribeSnapshotsOutput, Errs)
	GetElasticCacheCluster(input *elasticache.DescribeCacheClustersInput) ([]*elasticache.DescribeCacheClustersOutput, Errs)
	GetElasticacheTags(input *elasticache.ListTagsForResourceInput) ([]*elasticache.TagListMessage, Errs)
	GetLoadBalancers(input *elb.DescribeLoadBalancersInput) ([]*elb.DescribeLoadBalancersOutput, Errs)
	GetLoadBalancersTags(input *elb.DescribeTagsInput) ([]*elb.DescribeTagsOutput, Errs)
	GetLoadBalancersV2(input *elbv2.DescribeLoadBalancersInput) ([]*elbv2.DescribeLoadBalancersOutput, Errs)
	GetLoadBalancersV2Tags(input *elbv2.DescribeTagsInput) ([]*elbv2.DescribeTagsOutput, Errs)
	GetDBInstances(input *rds.DescribeDBInstancesInput) ([]*rds.DescribeDBInstancesOutput, Errs)
	GetDBInstancesTags(input *rds.ListTagsForResourceInput) ([]*rds.ListTagsForResourceOutput, Errs)
	ListBuckets(input *s3.ListBucketsInput) ([]*s3.ListBucketsOutput, Errs)
	GetBucketTags(input *s3.GetBucketTaggingInput) ([]*s3.GetBucketTaggingOutput, Errs)
	ListObjects(input *s3.ListObjectsInput) ([]*s3.ListObjectsOutput, Errs)
	DownloadObject(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error)
	GetObjectsTags(input *s3.GetObjectTaggingInput) ([]*s3.GetObjectTaggingOutput, Errs)
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
func NewAWSReader(accessKey string, secretKey string, regions []string, config *aws.Config) (AWSReader, error) {
	var c connector = connector{}

	creds, ec2s, sts, err := configureAWS(accessKey, secretKey)
	if err != nil {
		return nil, err
	}
	c.creds = creds
	if err := c.setAccountID(sts); err != nil {
		return nil, err
	}
	if err := c.setRegions(ec2s, regions); err != nil {
		return nil, err
	}
	c.setServices(config)
	return &c, nil
}

// GetAccountID returns the current ID for the account used
func (c *connector) GetAccountID() string {
	return *c.accountID
}

// GetRegions return the currently used regions for the Connector
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
	var token string = ""

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

func (c *connector) setRegions(ec2 ec2iface.EC2API, enabledRegions []string) error {
	if len(enabledRegions) == 0 {
		return errors.New("at least one region name is required")
	}
	regions, err := ec2.DescribeRegions(nil)
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

func (c *connector) setAccountID(sts stsiface.STSAPI) error {
	resp, err := sts.GetCallerIdentity(nil)
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
