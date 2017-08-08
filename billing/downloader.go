package billing

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cycloidio/raws"
)

type Downloader interface {
	Download(dest string) (string, error)
}

type billingDownloader struct {
	connector    raws.AWSReader
	s3Bucket     string
	filename     string
	fileFullPath string
}

func NewDownloader(s3Connector raws.AWSReader, bucket, filename string) Downloader {
	return &billingDownloader{
		connector: s3Connector,
		filename:  filename,
		s3Bucket:  bucket,
	}
}

func (d *billingDownloader) Download(dest string) (string, error) {
	fullPath, err := d.getAndCreateOutputPath(dest)
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
		Bucket: aws.String(d.s3Bucket),
		Key:    aws.String(d.filename)}
	_, err = d.connector.DownloadObject(fd, s3input)
	if err != nil {
		return "", fmt.Errorf("Error while downloading file: %+v", err)
	}
	return d.fileFullPath, nil
}

func Unzip(src string, dest string) (string, error) {
	unzipFullPath := dest + filepath.Base(strings.Split(src, ".zip")[0])
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := r.Close(); closeErr != nil {
			panic(closeErr)
		}
	}()

	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return "", err
	}

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if closeErr := f.Close(); closeErr != nil {
					panic(closeErr)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return "", err
		}
	}

	return unzipFullPath, nil
}

func (d *billingDownloader) getAndCreateOutputPath(dest string) (string, error) {
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
				return dest + d.filename, nil
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
		return dest + d.filename, nil
	}
	return dest, nil
}
