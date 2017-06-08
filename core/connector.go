package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/sts"
	"path/filepath"
)

type Connector struct {
	regions   []string
	svcs      []*serviceConnector
	creds     *credentials.Credentials
	accountID *string
}

func NewConnector(accessKey string, secretKey string, regions []string, config *aws.Config) (*Connector, error) {
	/* The default region is only used to (1) get the list of region and
	 * (2) get the account ID associated with the credentials.
	 *
	 * It is not used as a default region for services, therefore if no
	 * region is specified when instantiating the connector, then it will
	 * not try to establish any connections with AWS services.
	 */
	var defaultRegion string = "eu-west-1"
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
	region  string
	session *session.Session
	ec2     *ec2.EC2
	elb     *elb.ELB
	elbv2   *elbv2.ELBV2
}

func (c *Connector) setRegions(s *session.Session, enabledRegions []string) error {
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
