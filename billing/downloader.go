package billing

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"archive/zip"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cycloidio/raws"
)

type Downloader interface {
	Download(bucket string, filename string, dest string) (string, error)
	Unzip(src string, dest string) (string, error)
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

func (d *billingDownloader) Unzip(src string, dest string) (string, error) {
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
		rc, fileErr := f.Open()
		if fileErr != nil {
			return fileErr
		}
		defer func() {
			if closeErr := rc.Close(); closeErr != nil {
				panic(closeErr)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, openErr := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if openErr != nil {
				return openErr
			}
			defer func() {
				if closeErr := f.Close(); closeErr != nil {
					panic(closeErr)
				}
			}()

			_, copyErr := io.Copy(f, rc)
			if copyErr != nil {
				return copyErr
			}
		}
		return nil
	}

	for _, f := range r.File {
		writeErr := extractAndWriteFile(f)
		if writeErr != nil {
			return "", writeErr
		}
	}

	fi, err := os.Stat(unzipFullPath)
	if err != nil {
		return dest, nil
	}
	mode := fi.Mode()
	if !mode.IsDir() {
		return unzipFullPath, nil
	}
	return dest, nil
}
