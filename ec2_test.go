package raws

import (
	"errors"
	"testing"

	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type mockEC2 struct {
	ec2iface.EC2API

	// Mock of DescribeInstances
	dio   *ec2.DescribeInstancesOutput
	dierr error

	// Mock of DescribeVpcs
	dvpco   *ec2.DescribeVpcsOutput
	dvpcerr error

	// Mock of DescribeImages
	dimo   *ec2.DescribeImagesOutput
	dimerr error

	// Mocking of DescribeRegions
	dro   *ec2.DescribeRegionsOutput
	drerr error

	// Mock of DescribeSecurityGroups
	dsgo   *ec2.DescribeSecurityGroupsOutput
	dsgerr error

	// Mock DescribeSubnets
	dso   *ec2.DescribeSubnetsOutput
	dserr error

	// Mock DescribeVolumes
	dvolo   *ec2.DescribeVolumesOutput
	dvolerr error

	// Mock of DescribeSnapshots
	dsnapo   *ec2.DescribeSnapshotsOutput
	dsnaperr error
}

func (m mockEC2) DescribeImages(*ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
	return m.dimo, m.dimerr
}

func (m mockEC2) DescribeInstances(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return m.dio, m.dierr
}

// TODO: #17 -  Delete this mock after all the refactoring be done
func (m mockEC2) DescribeRegions(input *ec2.DescribeRegionsInput) (*ec2.DescribeRegionsOutput, error) {
	return m.dro, m.drerr
}

func (m mockEC2) DescribeRegionsWithContext(_ aws.Context, _ *ec2.DescribeRegionsInput, _ ...request.Option) (*ec2.DescribeRegionsOutput, error) {
	return m.dro, m.drerr
}

func (m mockEC2) DescribeSecurityGroups(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	return m.dsgo, m.dsgerr
}

func (m mockEC2) DescribeSnapshots(*ec2.DescribeSnapshotsInput) (*ec2.DescribeSnapshotsOutput, error) {
	return m.dsnapo, m.dsnaperr
}

func (m mockEC2) DescribeSubnets(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	return m.dso, m.dserr
}

func (m mockEC2) DescribeVolumes(*ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
	return m.dvolo, m.dvolerr
}

func (m mockEC2) DescribeVpcs(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	return m.dvpco, m.dvpcerr
}

func TestGetInstances(t *testing.T) {
	tests := []struct {
		name              string
		mocked            []*serviceConnector
		expectedInstances []*ec2.DescribeInstancesOutput
		expectedError     Errs
	}{{name: "one region with error",
		mocked: []*serviceConnector{
			{
				region: "test",
				ec2: mockEC2{
					dio:   &ec2.DescribeInstancesOutput{},
					dierr: errors.New("error with test"),
				},
			},
		},
		expectedError: Errs{&callErr{
			err:     errors.New("error with test"),
			region:  "test",
			service: ec2.ServiceName,
		}},
		expectedInstances: nil,
	},
		{name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					ec2: mockEC2{
						dio: &ec2.DescribeInstancesOutput{
							Reservations: []*ec2.Reservation{{
								OwnerId: aws.String("xxx"),
							}},
						},
						dierr: nil,
					},
				},
			},
			expectedError: nil,
			expectedInstances: []*ec2.DescribeInstancesOutput{
				{
					Reservations: []*ec2.Reservation{{
						OwnerId: aws.String("xxx"),
					}},
				},
			},
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					ec2: mockEC2{
						dio: &ec2.DescribeInstancesOutput{
							Reservations: []*ec2.Reservation{{
								OwnerId: aws.String("xxx"),
							}},
						},
						dierr: nil,
					},
				},
				{
					region: "test-2",
					ec2: mockEC2{
						dio: &ec2.DescribeInstancesOutput{
							Reservations: []*ec2.Reservation{{
								OwnerId: aws.String("yyy"),
							}},
						},
						dierr: nil,
					},
				},
			},
			expectedError: nil,
			expectedInstances: []*ec2.DescribeInstancesOutput{
				{
					Reservations: []*ec2.Reservation{{
						OwnerId: aws.String("xxx"),
					}},
				},
				{
					Reservations: []*ec2.Reservation{{
						OwnerId: aws.String("yyy"),
					}},
				},
			},
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					ec2: mockEC2{
						dio:   &ec2.DescribeInstancesOutput{},
						dierr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					ec2: mockEC2{
						dio: &ec2.DescribeInstancesOutput{
							Reservations: []*ec2.Reservation{{
								OwnerId: aws.String("yyy"),
							}},
						},
						dierr: nil,
					},
				},
			},
			expectedError: Errs{&callErr{
				err:     errors.New("error with test-1"),
				region:  "test-1",
				service: ec2.ServiceName,
			}},
			expectedInstances: []*ec2.DescribeInstancesOutput{
				{
					Reservations: []*ec2.Reservation{{
						OwnerId: aws.String("yyy"),
					}},
				},
			},
		}}

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		instances, err := c.GetInstances(nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(instances, tt.expectedInstances) {
			t.Errorf("%s [%d] - EC2 instances: received=%+v | expected=%+v",
				tt.name, i, instances, tt.expectedInstances)
		}
	}
}

func TestGetVpcs(t *testing.T) {
	tests := []struct {
		name          string
		mocked        []*serviceConnector
		expectedVpcs  []*ec2.DescribeVpcsOutput
		expectedError Errs
	}{{name: "one region with error",
		mocked: []*serviceConnector{
			{
				region: "test",
				ec2: mockEC2{
					dvpco:   &ec2.DescribeVpcsOutput{},
					dvpcerr: errors.New("error with test"),
				},
			},
		},
		expectedError: Errs{&callErr{
			err:     errors.New("error with test"),
			region:  "test",
			service: ec2.ServiceName,
		}},
		expectedVpcs: nil,
	},
		{name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					ec2: mockEC2{
						dvpco: &ec2.DescribeVpcsOutput{
							Vpcs: []*ec2.Vpc{{
								VpcId: aws.String("1"),
							}},
						},
						dvpcerr: nil,
					},
				},
			},
			expectedError: nil,
			expectedVpcs: []*ec2.DescribeVpcsOutput{
				{
					Vpcs: []*ec2.Vpc{{
						VpcId: aws.String("1"),
					}},
				},
			},
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					ec2: mockEC2{
						dvpco: &ec2.DescribeVpcsOutput{
							Vpcs: []*ec2.Vpc{{
								VpcId: aws.String("1"),
							}},
						},
						dvpcerr: nil,
					},
				},
				{
					region: "test-2",
					ec2: mockEC2{
						dvpco: &ec2.DescribeVpcsOutput{
							Vpcs: []*ec2.Vpc{{
								VpcId: aws.String("2"),
							}},
						},
						dvpcerr: nil,
					},
				},
			},
			expectedError: nil,
			expectedVpcs: []*ec2.DescribeVpcsOutput{
				{
					Vpcs: []*ec2.Vpc{{
						VpcId: aws.String("1"),
					}},
				},
				{
					Vpcs: []*ec2.Vpc{{
						VpcId: aws.String("2"),
					}},
				},
			},
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					ec2: mockEC2{
						dvpco:   &ec2.DescribeVpcsOutput{},
						dvpcerr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					ec2: mockEC2{
						dvpco: &ec2.DescribeVpcsOutput{
							Vpcs: []*ec2.Vpc{{
								VpcId: aws.String("2"),
							}},
						},
						dvpcerr: nil,
					},
				},
			},
			expectedError: Errs{&callErr{
				err:     errors.New("error with test-1"),
				region:  "test-1",
				service: ec2.ServiceName,
			}},
			expectedVpcs: []*ec2.DescribeVpcsOutput{
				{
					Vpcs: []*ec2.Vpc{{
						VpcId: aws.String("2"),
					}},
				},
			},
		}}

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		vpcs, err := c.GetVpcs(nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(vpcs, tt.expectedVpcs) {
			t.Errorf("%s [%d] - EC2 VPCs: received=%+v | expected=%+v",
				tt.name, i, vpcs, tt.expectedVpcs)
		}
	}
}

func TestGetImages(t *testing.T) {
	tests := []struct {
		name           string
		mocked         []*serviceConnector
		expectedImages []*ec2.DescribeImagesOutput
		expectedError  Errs
	}{{name: "one region with error",
		mocked: []*serviceConnector{
			{
				region: "test",
				ec2: mockEC2{
					dimo:   &ec2.DescribeImagesOutput{},
					dimerr: errors.New("error with test"),
				},
			},
		},
		expectedError: Errs{&callErr{
			err:     errors.New("error with test"),
			region:  "test",
			service: ec2.ServiceName,
		}},
		expectedImages: nil,
	},
		{name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					ec2: mockEC2{
						dimo: &ec2.DescribeImagesOutput{
							Images: []*ec2.Image{{
								Name: aws.String("test"),
							}},
						},
						dvpcerr: nil,
					},
				},
			},
			expectedError: nil,
			expectedImages: []*ec2.DescribeImagesOutput{
				{
					Images: []*ec2.Image{{
						Name: aws.String("test"),
					}},
				},
			},
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					ec2: mockEC2{
						dimo: &ec2.DescribeImagesOutput{
							Images: []*ec2.Image{{
								Name: aws.String("test-1"),
							}},
						},
						dimerr: nil,
					},
				},
				{
					region: "test-2",
					ec2: mockEC2{
						dimo: &ec2.DescribeImagesOutput{
							Images: []*ec2.Image{{
								Name: aws.String("test-2"),
							}},
						},
						dimerr: nil,
					},
				},
			},
			expectedError: nil,
			expectedImages: []*ec2.DescribeImagesOutput{
				{
					Images: []*ec2.Image{{
						Name: aws.String("test-1"),
					}},
				},
				{
					Images: []*ec2.Image{{
						Name: aws.String("test-2"),
					}},
				},
			},
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					ec2: mockEC2{
						dimo:   &ec2.DescribeImagesOutput{},
						dimerr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					ec2: mockEC2{
						dimo: &ec2.DescribeImagesOutput{
							Images: []*ec2.Image{{
								Name: aws.String("test-2"),
							}},
						},
						dimerr: nil,
					},
				},
			},
			expectedError: Errs{&callErr{
				err:     errors.New("error with test-1"),
				region:  "test-1",
				service: ec2.ServiceName,
			}},
			expectedImages: []*ec2.DescribeImagesOutput{
				{
					Images: []*ec2.Image{{
						Name: aws.String("test-2"),
					}},
				},
			},
		}}

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		images, err := c.GetImages(nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(images, tt.expectedImages) {
			t.Errorf("%s [%d] - EC2 Images: received=%+v | expected=%+v",
				tt.name, i, images, tt.expectedImages)
		}
	}
}

func TestGetSecurityGroups(t *testing.T) {
	tests := []struct {
		name              string
		mocked            []*serviceConnector
		expectedSecGroups []*ec2.DescribeSecurityGroupsOutput
		expectedError     Errs
	}{{name: "one region with error",
		mocked: []*serviceConnector{
			{
				region: "test",
				ec2: mockEC2{
					dsgo:   &ec2.DescribeSecurityGroupsOutput{},
					dsgerr: errors.New("error with test"),
				},
			},
		},
		expectedError: Errs{&callErr{
			err:     errors.New("error with test"),
			region:  "test",
			service: ec2.ServiceName,
		}},
		expectedSecGroups: nil,
	},
		{name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					ec2: mockEC2{
						dsgo: &ec2.DescribeSecurityGroupsOutput{
							SecurityGroups: []*ec2.SecurityGroup{{
								GroupId: aws.String("1"),
							}},
						},
						dsgerr: nil,
					},
				},
			},
			expectedError: nil,
			expectedSecGroups: []*ec2.DescribeSecurityGroupsOutput{
				{
					SecurityGroups: []*ec2.SecurityGroup{{
						GroupId: aws.String("1"),
					}},
				},
			},
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					ec2: mockEC2{
						dsgo: &ec2.DescribeSecurityGroupsOutput{
							SecurityGroups: []*ec2.SecurityGroup{{
								GroupId: aws.String("1"),
							}},
						},
						dsgerr: nil,
					},
				},
				{
					region: "test-2",
					ec2: mockEC2{
						dsgo: &ec2.DescribeSecurityGroupsOutput{
							SecurityGroups: []*ec2.SecurityGroup{{
								GroupId: aws.String("2"),
							}},
						},
						dsgerr: nil,
					},
				},
			},
			expectedError: nil,
			expectedSecGroups: []*ec2.DescribeSecurityGroupsOutput{
				{
					SecurityGroups: []*ec2.SecurityGroup{{
						GroupId: aws.String("1"),
					}},
				},
				{
					SecurityGroups: []*ec2.SecurityGroup{{
						GroupId: aws.String("2"),
					}},
				},
			},
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					ec2: mockEC2{
						dsgo:   &ec2.DescribeSecurityGroupsOutput{},
						dsgerr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					ec2: mockEC2{
						dsgo: &ec2.DescribeSecurityGroupsOutput{
							SecurityGroups: []*ec2.SecurityGroup{{
								GroupId: aws.String("2"),
							}},
						},
						dsgerr: nil,
					},
				},
			},
			expectedError: Errs{&callErr{
				err:     errors.New("error with test-1"),
				region:  "test-1",
				service: ec2.ServiceName,
			}},
			expectedSecGroups: []*ec2.DescribeSecurityGroupsOutput{
				{
					SecurityGroups: []*ec2.SecurityGroup{{
						GroupId: aws.String("2"),
					}},
				},
			},
		}}

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		secGroups, err := c.GetSecurityGroups(nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(secGroups, tt.expectedSecGroups) {
			t.Errorf("%s [%d] - EC2 security groups: received=%+v | expected=%+v",
				tt.name, i, secGroups, tt.expectedSecGroups)
		}
	}
}

func TestGetSubnets(t *testing.T) {
	tests := []struct {
		name            string
		mocked          []*serviceConnector
		expectedSubnets []*ec2.DescribeSubnetsOutput
		expectedError   Errs
	}{{name: "one region with error",
		mocked: []*serviceConnector{
			{
				region: "test",
				ec2: mockEC2{
					dso:   &ec2.DescribeSubnetsOutput{},
					dserr: errors.New("error with test"),
				},
			},
		},
		expectedError: Errs{&callErr{
			err:     errors.New("error with test"),
			region:  "test",
			service: ec2.ServiceName,
		}},
		expectedSubnets: nil,
	},
		{name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					ec2: mockEC2{
						dso: &ec2.DescribeSubnetsOutput{
							Subnets: []*ec2.Subnet{{
								SubnetId: aws.String("1"),
							}},
						},
						dserr: nil,
					},
				},
			},
			expectedError: nil,
			expectedSubnets: []*ec2.DescribeSubnetsOutput{
				{
					Subnets: []*ec2.Subnet{{
						SubnetId: aws.String("1"),
					}},
				},
			},
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					ec2: mockEC2{
						dso: &ec2.DescribeSubnetsOutput{
							Subnets: []*ec2.Subnet{{
								SubnetId: aws.String("1"),
							}},
						},
						dserr: nil,
					},
				},
				{
					region: "test-2",
					ec2: mockEC2{
						dso: &ec2.DescribeSubnetsOutput{
							Subnets: []*ec2.Subnet{{
								SubnetId: aws.String("2"),
							}},
						},
						dserr: nil,
					},
				},
			},
			expectedError: nil,
			expectedSubnets: []*ec2.DescribeSubnetsOutput{
				{
					Subnets: []*ec2.Subnet{{
						SubnetId: aws.String("1"),
					}},
				},
				{
					Subnets: []*ec2.Subnet{{
						SubnetId: aws.String("2"),
					}},
				},
			},
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					ec2: mockEC2{
						dso:   &ec2.DescribeSubnetsOutput{},
						dserr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					ec2: mockEC2{
						dso: &ec2.DescribeSubnetsOutput{
							Subnets: []*ec2.Subnet{{
								SubnetId: aws.String("2"),
							}},
						},
						dserr: nil,
					},
				},
			},
			expectedError: Errs{&callErr{
				err:     errors.New("error with test-1"),
				region:  "test-1",
				service: ec2.ServiceName,
			}},
			expectedSubnets: []*ec2.DescribeSubnetsOutput{
				{
					Subnets: []*ec2.Subnet{{
						SubnetId: aws.String("2"),
					}},
				},
			},
		}}

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		subnets, err := c.GetSubnets(nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(subnets, tt.expectedSubnets) {
			t.Errorf("%s [%d] - EC2 subnets: received=%+v | expected=%+v",
				tt.name, i, subnets, tt.expectedSubnets)
		}
	}
}

func TestGetVolumes(t *testing.T) {
	tests := []struct {
		name            string
		mocked          []*serviceConnector
		expectedVolumes []*ec2.DescribeVolumesOutput
		expectedError   Errs
	}{{name: "one region with error",
		mocked: []*serviceConnector{
			{
				region: "test",
				ec2: mockEC2{
					dvolo:   &ec2.DescribeVolumesOutput{},
					dvolerr: errors.New("error with test"),
				},
			},
		},
		expectedError: Errs{&callErr{
			err:     errors.New("error with test"),
			region:  "test",
			service: ec2.ServiceName,
		}},
		expectedVolumes: nil,
	},
		{name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					ec2: mockEC2{
						dvolo: &ec2.DescribeVolumesOutput{
							Volumes: []*ec2.Volume{{
								VolumeId: aws.String("1"),
							}},
						},
						dvolerr: nil,
					},
				},
			},
			expectedError: nil,
			expectedVolumes: []*ec2.DescribeVolumesOutput{
				{
					Volumes: []*ec2.Volume{{
						VolumeId: aws.String("1"),
					}},
				},
			},
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					ec2: mockEC2{
						dvolo: &ec2.DescribeVolumesOutput{
							Volumes: []*ec2.Volume{{
								VolumeId: aws.String("1"),
							}},
						},
						dvolerr: nil,
					},
				},
				{
					region: "test-2",
					ec2: mockEC2{
						dvolo: &ec2.DescribeVolumesOutput{
							Volumes: []*ec2.Volume{{
								VolumeId: aws.String("2"),
							}},
						},
						dvolerr: nil,
					},
				},
			},
			expectedError: nil,
			expectedVolumes: []*ec2.DescribeVolumesOutput{
				{
					Volumes: []*ec2.Volume{{
						VolumeId: aws.String("1"),
					}},
				},
				{
					Volumes: []*ec2.Volume{{
						VolumeId: aws.String("2"),
					}},
				},
			},
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					ec2: mockEC2{
						dvolo:   &ec2.DescribeVolumesOutput{},
						dvolerr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					ec2: mockEC2{
						dvolo: &ec2.DescribeVolumesOutput{
							Volumes: []*ec2.Volume{{
								VolumeId: aws.String("2"),
							}},
						},
						dvolerr: nil,
					},
				},
			},
			expectedError: Errs{&callErr{
				err:     errors.New("error with test-1"),
				region:  "test-1",
				service: ec2.ServiceName,
			}},
			expectedVolumes: []*ec2.DescribeVolumesOutput{
				{
					Volumes: []*ec2.Volume{{
						VolumeId: aws.String("2"),
					}},
				},
			},
		}}

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		volumes, err := c.GetVolumes(nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(volumes, tt.expectedVolumes) {
			t.Errorf("%s [%d] - EC2 volumes: received=%+v | expected=%+v",
				tt.name, i, volumes, tt.expectedVolumes)
		}
	}
}

func TestGetSnapshots(t *testing.T) {
	tests := []struct {
		name              string
		mocked            []*serviceConnector
		expectedSnapshots []*ec2.DescribeSnapshotsOutput
		expectedError     Errs
	}{{name: "one region with error",
		mocked: []*serviceConnector{
			{
				region: "test",
				ec2: mockEC2{
					dsnapo: &ec2.DescribeSnapshotsOutput{
						Snapshots: []*ec2.Snapshot{{
							SnapshotId: aws.String("1"),
						}},
					},
					dsnaperr: errors.New("error with test"),
				},
			},
		},
		expectedError: Errs{&callErr{
			err:     errors.New("error with test"),
			region:  "test",
			service: ec2.ServiceName,
		}},
		expectedSnapshots: nil,
	},
		{name: "one region no error",
			mocked: []*serviceConnector{
				{
					region: "test",
					ec2: mockEC2{
						dsnapo: &ec2.DescribeSnapshotsOutput{
							Snapshots: []*ec2.Snapshot{{
								SnapshotId: aws.String("1"),
							}},
						},
						dsnaperr: nil,
					},
				},
			},
			expectedError: nil,
			expectedSnapshots: []*ec2.DescribeSnapshotsOutput{
				{
					Snapshots: []*ec2.Snapshot{{
						SnapshotId: aws.String("1"),
					}},
				},
			},
		},
		{name: "multiple region no error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					ec2: mockEC2{
						dsnapo: &ec2.DescribeSnapshotsOutput{
							Snapshots: []*ec2.Snapshot{{
								SnapshotId: aws.String("1"),
							}},
						},
						dsnaperr: nil,
					},
				},
				{
					region: "test-2",
					ec2: mockEC2{
						dsnapo: &ec2.DescribeSnapshotsOutput{
							Snapshots: []*ec2.Snapshot{{
								SnapshotId: aws.String("2"),
							}},
						},
						dsnaperr: nil,
					},
				},
			},
			expectedError: nil,
			expectedSnapshots: []*ec2.DescribeSnapshotsOutput{
				{
					Snapshots: []*ec2.Snapshot{{
						SnapshotId: aws.String("1"),
					}},
				},
				{
					Snapshots: []*ec2.Snapshot{{
						SnapshotId: aws.String("2"),
					}},
				},
			},
		},
		{name: "multiple region with error",
			mocked: []*serviceConnector{
				{
					region: "test-1",
					ec2: mockEC2{
						dsnapo:   &ec2.DescribeSnapshotsOutput{},
						dsnaperr: errors.New("error with test-1"),
					},
				},
				{
					region: "test-2",
					ec2: mockEC2{
						dsnapo: &ec2.DescribeSnapshotsOutput{
							Snapshots: []*ec2.Snapshot{{
								SnapshotId: aws.String("2"),
							}},
						},
						dsnaperr: nil,
					},
				},
			},
			expectedError: Errs{&callErr{
				err:     errors.New("error with test-1"),
				region:  "test-1",
				service: ec2.ServiceName,
			}},
			expectedSnapshots: []*ec2.DescribeSnapshotsOutput{
				{
					Snapshots: []*ec2.Snapshot{{
						SnapshotId: aws.String("2"),
					}},
				},
			},
		}}

	for i, tt := range tests {
		c := &connector{svcs: tt.mocked}
		snapshots, err := c.GetSnapshots(nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if !reflect.DeepEqual(snapshots, tt.expectedSnapshots) {
			t.Errorf("%s [%d] - EC2 snapshots: received=%+v | expected=%+v",
				tt.name, i, snapshots, tt.expectedSnapshots)
		}
	}
}
