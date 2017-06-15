package core

import (
	"path/filepath"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elasticache/elasticacheiface"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
)

// The Connector provides easy access to AWS SDK calls.
//
// By using it, calls can be made directly through multiple regions, and will filter only data that belongs to you.
// For example, when fetching the list of AMI, or snapshots.
//
// In order to start making calls, only calling NewConnector is required.
type Connector struct {
	regions   []string
	svcs      []*serviceConnector
	creds     *credentials.Credentials
	accountID *string
}

// NewConnector returns an object which also contains the accountID and extend the different regions to use.
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
func NewConnector(accessKey string, secretKey string, regions []string, config *aws.Config) (*Connector, error) {
	/* The default region is only used to (1) get the list of region and
	 * (2) get the account ID associated with the credentials.
	 *
	 * It is not used as a default region for services, therefore if no
	 * region is specified when instantiating the connector, then it will
	 * not try to establish any connections with AWS services.
	 */
	const defaultRegion string = "eu-west-1"
	var defaultSession *session.Session
	var c Connector = Connector{}

	c.setCredentials(accessKey, secretKey)
	defaultSession = session.Must(
		session.NewSession(&aws.Config{
			Region:      aws.String(defaultRegion),
			DisableSSL:  aws.Bool(false),
			MaxRetries:  aws.Int(3),
			Credentials: c.creds,
		}),
	)
	if err := c.setAccountID(defaultSession); err != nil {
		return nil, err
	}
	if err := c.setRegions(defaultSession, regions); err != nil {
		return nil, err
	}
	c.setServices(config)
	return &c, nil
}

type serviceConnector struct {
	region      string
	session     *session.Session
	ec2         ec2iface.EC2API
	elb         elbiface.ELBAPI
	elbv2       elbv2iface.ELBV2API
	rds         rdsiface.RDSAPI
	s3          s3iface.S3API
	elasticache elasticacheiface.ElastiCacheAPI
}

func (c *Connector) setRegions(s *session.Session, enabledRegions []string) error {
	if len(enabledRegions) == 0 {
		return errors.New("at least one region name is required")
	}
	svc := ec2.New(s)
	regions, err := svc.DescribeRegions(nil)

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

func (c *Connector) setAccountID(s *session.Session) error {
	var params *sts.GetCallerIdentityInput

	svc := sts.New(s)
	resp, err := svc.GetCallerIdentity(params)
	if err != nil {
		return err
	}
	c.accountID = resp.Account
	return nil
}

func (c *Connector) setCredentials(accessKey string, secretKey string) error {
	var token string = ""

	c.creds = credentials.NewStaticCredentials(accessKey, secretKey, token)

	_, err := c.creds.Get()
	if err != nil {
		return err
	}
	return nil
}

func (c *Connector) setServices(config *aws.Config) error {
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
		session := session.Must(session.NewSession(config))
		svc := &serviceConnector{
			region:  region,
			session: session,
		}
		c.svcs = append(c.svcs, svc)
	}
	return nil
}
