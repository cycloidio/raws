package raws

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (c *connector) DownloadObject(
	ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader),
) (int64, error) {
	var err error
	var n int64

	if input.Bucket == nil || input.Key == nil {
		return n, fmt.Errorf("couldn't download undefined object (keys or bucket not set)")
	}
	for _, svc := range c.svcs {
		if svc.s3downloader == nil {
			svc.s3downloader = s3manager.NewDownloader(svc.session)
		}
		n, err = svc.s3downloader.DownloadWithContext(ctx, w, input, options...)
		if err == nil {
			return n, nil
		}
	}
	return n, fmt.Errorf("couldn't download '%s/%s' in any of '%+v' regions", *input.Bucket, *input.Key, c.GetRegions())
}
