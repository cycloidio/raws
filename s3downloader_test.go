package raws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

type mockS3Downloader struct {
	s3manageriface.DownloaderAPI

	// Mock of DownloadObject
	dob   int64
	doerr error
}

func (m mockS3Downloader) DownloadWithContext(
	_ aws.Context, _ io.WriterAt, _ *s3.GetObjectInput, _ ...func(*s3manager.Downloader),
) (int64, error) {
	return m.dob, m.doerr
}

func TestDownloadObjet(t *testing.T) {

	tests := []struct {
		name          string
		mocked        []*serviceConnector
		regions       []string
		input         *s3.GetObjectInput
		expectedBytes int64
		expectedError error
	}{
		{
			name: "one region with error",
			mocked: []*serviceConnector{
				{
					s3downloader: mockS3Downloader{
						dob:   0,
						doerr: errors.New("error with test"),
					},
				},
			},
			regions: []string{"test"},
			input: &s3.GetObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
			},
			expectedError: fmt.Errorf("couldn't download 'bucket/key' in any of '[test]' regions"),
			expectedBytes: 0,
		},
		{
			name:    "invalid file requested",
			mocked:  nil,
			regions: []string{"test"},
			input: &s3.GetObjectInput{
				Bucket: nil,
				Key:    nil,
			},
			expectedError: fmt.Errorf("couldn't download undefined object (keys or bucket not set)"),
			expectedBytes: 0,
		},
		{
			name: "one region no error",
			mocked: []*serviceConnector{
				{
					s3downloader: mockS3Downloader{
						dob:   42,
						doerr: nil,
					},
				},
			},
			regions: []string{"test"},
			input: &s3.GetObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
			},
			expectedError: nil,
			expectedBytes: 42,
		},
	}

	var ctx = context.Background()

	for i, tt := range tests {
		c := &connector{
			regions: tt.regions,
			svcs:    tt.mocked,
		}
		bytes, err := c.DownloadObject(ctx, nil, tt.input, nil)
		checkErrors(t, tt.name, i, err, tt.expectedError)
		if tt.expectedBytes != bytes {
			t.Errorf("%s [%d] - S3 download object: received=%+v | expected=%+v",
				tt.name, i, bytes, tt.expectedBytes)
		}
	}
}
