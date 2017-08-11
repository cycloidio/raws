package billing

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cycloidio/raws"
)

type Downloader interface {
	Download(bucket string, filename string, dest string) (string, error)
}

type billingDownloader struct {
	connector    raws.AWSReader
	fileFullPath string
}

func NewDownloader(s3Connector raws.AWSReader) Downloader {
	return &billingDownloader{
		connector: s3Connector,
	}
}

func (d *billingDownloader) Download(bucket string, filename string, dest string) (string, error) {
	fullPath, err := d.getAndCreateOutputPath(filename, dest)
	if err != nil {
		return "", fmt.Errorf("Error while identifying destination's path: %v", err)
	}
	d.fileFullPath = fullPath
	fd, err := os.Create(d.fileFullPath)
	if err != nil {
		return "", fmt.Errorf("Couldn't create file %q: %+v", dest, err)
	}
	defer fd.Close()
	s3input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename)}
	_, err = d.connector.DownloadObject(fd, s3input)
	if err != nil {
		return "", fmt.Errorf("Error while downloading file: %+v", err)
	}
	return d.fileFullPath, nil
}

func (d *billingDownloader) getAndCreateOutputPath(filename string, dest string) (string, error) {
	fi, err := os.Stat(dest)
	if err != nil {
		// Case when the destination doesn't exist
		// The directory tree is first created, and then the path is
		// checked, in order to know if we were given a supposed 'file'
		// or if we were given a directory.
		if os.IsNotExist(err) {
			osErr := os.MkdirAll(filepath.Dir(dest), 0755)
			if osErr != nil {
				return "", osErr
			}
			if strings.Contains(filepath.Dir(dest), filepath.Base(dest)) {
				return dest + filename, nil
			} else {
				return dest, nil
			}
		}
		return "", err
	}
	// Case when the destination does exist
	// If the path is a directory then we return it with the default path,
	// otherwise if it is a file, we simply returns it.
	mode := fi.Mode()
	if mode.IsDir() {
		return dest + filename, nil
	}
	return dest, nil
}
