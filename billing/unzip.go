package billing

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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
