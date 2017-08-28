package billing

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cycloidio/raws"
)

func TestBillingChecker_Check(t *testing.T) {
	const (
		givenFilename = "test-file"
		givenBucket   = "test-bucket"
	)
	tests := []struct {
		name           string
		mockedDyn      mockDynamodb
		mockedS3       mockAWSReader
		expectedCheck  bool
		expectedError  error
		expectedS3MD5  string
		expectedDynMD5 string
	}{
		{name: "no error while getting md5(s) - different md5",
			mockedS3: mockAWSReader{
				loo: []*s3.ListObjectsOutput{
					{
						Contents: []*s3.Object{
							&s3.Object{
								ETag: aws.String("\"1111111111111\""),
							},
						},
					},
				},
				loe: nil,
			},
			mockedDyn: mockDynamodb{
				geo: &dynamodb.GetItemOutput{
					Item: map[string]*dynamodb.AttributeValue{
						billingReportMd5Field: &dynamodb.AttributeValue{
							S: aws.String("2222222222222"),
						},
					},
				},
				gee: nil,
			},
			expectedCheck:  true,
			expectedError:  nil,
			expectedS3MD5:  "1111111111111",
			expectedDynMD5: "2222222222222",
		},
		{name: "error while getting s3 md5",
			mockedS3: mockAWSReader{
				loo: nil,
				loe: raws.Errs{
					raws.NewAPIError("", s3.ServiceName, errors.New("error while getting objects")),
				},
			},
			mockedDyn: mockDynamodb{
				geo: &dynamodb.GetItemOutput{
					Item: map[string]*dynamodb.AttributeValue{
						billingReportMd5Field: &dynamodb.AttributeValue{
							S: aws.String("2222222222222"),
						},
					},
				},
				gee: nil,
			},
			expectedCheck: false,
			expectedError: raws.Errs{
				raws.NewAPIError("", s3.ServiceName, errors.New("error while getting objects")),
			},
			expectedS3MD5:  "",
			expectedDynMD5: "2222222222222",
		},
		{name: "error while getting dynamodb md5",
			// Will not reached that part
			mockedS3: mockAWSReader{},
			mockedDyn: mockDynamodb{
				geo: &dynamodb.GetItemOutput{
					Item: map[string]*dynamodb.AttributeValue{
						"key": &dynamodb.AttributeValue{S: aws.String("value")},
					},
				},
				gee: nil,
			},
			expectedCheck:  false,
			expectedError:  fmt.Errorf("no '%s' field present for the entity", billingReportMd5Field),
			expectedS3MD5:  "",
			expectedDynMD5: "",
		},
		{name: "no error while getting md5(s) - same md5",
			mockedS3: mockAWSReader{
				loo: []*s3.ListObjectsOutput{
					{
						Contents: []*s3.Object{
							&s3.Object{
								ETag: aws.String("\"1111111111111\""),
							},
						},
					},
				},
				loe: nil,
			},
			mockedDyn: mockDynamodb{
				geo: &dynamodb.GetItemOutput{
					Item: map[string]*dynamodb.AttributeValue{
						billingReportMd5Field: &dynamodb.AttributeValue{
							S: aws.String("1111111111111"),
						},
					},
				},
				gee: nil,
			},
			expectedCheck:  false,
			expectedError:  nil,
			expectedS3MD5:  "1111111111111",
			expectedDynMD5: "1111111111111",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &billingChecker{
				s3Connector: tt.mockedS3,
				dynSvc:      tt.mockedDyn,
			}
			needCheck, err := c.Check(givenBucket, givenFilename)
			checkErrors(t, tt.name, i, err, tt.expectedError)
			if needCheck != tt.expectedCheck {
				t.Errorf("%s [%d] - Check boolean doesn't match: received=%t | expected=%t",
					tt.name, i, needCheck, tt.expectedCheck)
			}
			if c.oldMd5 != tt.expectedDynMD5 {
				t.Errorf("%s [%d] - Check old md5 doesn't match: received=%q | expected=%q",
					tt.name, i, c.oldMd5, tt.expectedDynMD5)
			}
			if c.newMd5 != tt.expectedS3MD5 {
				t.Errorf("%s [%d] - Check new md5 doesn't match: received=%q | expected=%q",
					tt.name, i, c.newMd5, tt.expectedS3MD5)
			}
		})
	}
}

func TestBillingChecker_AlreadyPresent(t *testing.T) {
	tests := []struct {
		name             string
		oldMd5           string
		newMd5           string
		expectedPresence bool
		expectedMd5      string
	}{
		{name: "new & old md5 are identical",
			oldMd5:           "11111111111",
			newMd5:           "11111111111",
			expectedPresence: true,
			expectedMd5:      "11111111111",
		},
		{name: "new & old md5 are different",
			oldMd5:           "11111111111",
			newMd5:           "22222222222",
			expectedPresence: false,
			expectedMd5:      "22222222222",
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &billingChecker{
				oldMd5: tt.oldMd5,
				newMd5: tt.newMd5,
			}
			present, md5 := c.AlreadyPresent()
			if present != tt.expectedPresence {
				t.Errorf("%s [%d] - presence invalid: received=%t | expected=%t",
					tt.name, i, present, tt.expectedPresence)
			}
			if md5 != tt.expectedMd5 {
				t.Errorf("%s [%d] - md5 incorrect: received=%q | expected=%q",
					tt.name, i, md5, tt.expectedMd5)
			}
		})
	}
}

func TestBillingChecker_getS3Entry(t *testing.T) {
	const (
		givenFilename = "test-file"
		givenBucket   = "test-bucket"
	)
	tests := []struct {
		name          string
		mockedS3      mockAWSReader
		expectedMD5   string
		expectedError error
	}{
		{name: "no error while getting md5",
			mockedS3: mockAWSReader{
				loo: []*s3.ListObjectsOutput{
					{
						Contents: []*s3.Object{
							&s3.Object{
								ETag: aws.String("\"1111111111111\""),
							},
						},
					},
				},
				loe: nil,
			},
			expectedMD5:   "1111111111111",
			expectedError: nil,
		},
		{name: "error while getting list objects",
			mockedS3: mockAWSReader{
				loo: nil,
				loe: raws.Errs{
					raws.NewAPIError("", s3.ServiceName, errors.New("error while getting objects")),
				},
			},
			expectedMD5: "",
			expectedError: raws.Errs{
				raws.NewAPIError("", s3.ServiceName, errors.New("error while getting objects")),
			},
		},
		{name: "too many objects",
			mockedS3: mockAWSReader{
				loo: []*s3.ListObjectsOutput{{}, {}},
				loe: nil,
			},
			expectedMD5:   "",
			expectedError: errors.New("found too many objects matching (2)"),
		},
		{name: "content is nil",
			mockedS3: mockAWSReader{
				loo: []*s3.ListObjectsOutput{{}},
				loe: nil,
			},
			expectedMD5:   "",
			expectedError: errors.New("s3 entry doesn't have 'Contents' attribute"),
		}}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &billingChecker{
				s3Connector: tt.mockedS3,
			}
			err := c.getS3Entry(givenBucket, givenFilename)
			checkErrors(t, tt.name, i, err, tt.expectedError)
			if !reflect.DeepEqual(c.newMd5, tt.expectedMD5) {
				t.Errorf("%s [%d] - s3 md5: received=%q | expected=%q",
					tt.name, i, c.newMd5, tt.expectedMD5)
			}
		})
	}
}

func TestBillingChecker_getDynamoEntry(t *testing.T) {
	const (
		givenFilename = "test-file"
	)
	tests := []struct {
		name          string
		mockedDyn     mockDynamodb
		expectedMD5   string
		expectedError error
	}{
		{name: "no error while getting md5",
			mockedDyn: mockDynamodb{
				geo: &dynamodb.GetItemOutput{
					Item: map[string]*dynamodb.AttributeValue{
						billingReportMd5Field: &dynamodb.AttributeValue{
							S: aws.String("1111111111111"),
						},
					},
				},
				gee: nil,
			},
			expectedMD5:   "1111111111111",
			expectedError: nil,
		},
		{name: "error while getting list objects",
			mockedDyn: mockDynamodb{
				geo: nil,
				gee: errors.New("error while getting item"),
			},
			expectedMD5:   "",
			expectedError: errors.New("error while getting item"),
		},
		{name: "no single item return",
			mockedDyn: mockDynamodb{
				geo: &dynamodb.GetItemOutput{Item: nil},
				gee: nil,
			},
			expectedMD5:   "",
			expectedError: nil,
		},
		{name: "item returned doesn't have md5 field",
			mockedDyn: mockDynamodb{
				geo: &dynamodb.GetItemOutput{
					Item: map[string]*dynamodb.AttributeValue{
						"key": &dynamodb.AttributeValue{S: aws.String("value")},
					},
				},
				gee: nil,
			},
			expectedMD5:   "",
			expectedError: fmt.Errorf("no '%s' field present for the entity", billingReportMd5Field),
		}}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &billingChecker{
				dynSvc: tt.mockedDyn,
			}
			err := c.getDynamoEntry(givenFilename)
			checkErrors(t, tt.name, i, err, tt.expectedError)
			if !reflect.DeepEqual(c.oldMd5, tt.expectedMD5) {
				t.Errorf("%s [%d] - dynamodb md5: received=%q | expected=%q",
					tt.name, i, c.oldMd5, tt.expectedMD5)
			}
		})
	}
}
