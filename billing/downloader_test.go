package billing

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/cycloidio/raws"
)

func TestNewDownloader(t *testing.T) {
	var mockedS3 raws.AWSReader = mockAWSReader{}

	d := &billingDownloader{
		connector: mockedS3,
	}
	cd := NewDownloader(mockedS3)
	if !reflect.DeepEqual(d, cd) {
		t.Errorf("NewDownloader: received=%+v | expected=%+v",
			cd, d)
	}
}

func TestBillingDownloader_Download(t *testing.T) {
	const (
		tempDownloadDir string = "/billingDownloader/"
		tempFilename    string = "test.csv.zip"
		givenFilename   string = "bd.csv.zip"
		givenBucket     string = "fake-bucket"
	)
	var tempDir = os.TempDir() + tempDownloadDir
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		t.Errorf("Error while creating temporary dir: %q - %v", tempDir, err)
	}
	var tempFile = tempDir + tempFilename

	tests := []struct {
		name          string
		destination   string
		mockS3        mockAWSReader
		expectedPath  string
		expectedError error
	}{
		{name: "no errors during download",
			mockS3: mockAWSReader{
				doi: 0,
				doe: nil,
			},
			destination:   tempFile,
			expectedPath:  tempFile,
			expectedError: nil,
		},
		{name: "errors during download",
			mockS3: mockAWSReader{
				doi: 0,
				doe: errors.New("test error"),
			},
			destination:   tempFile,
			expectedPath:  "",
			expectedError: errors.New("Error while downloading file: test error"),
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &billingDownloader{
				connector: tt.mockS3,
			}
			path, err := d.Download(givenBucket, givenFilename, tt.destination)
			checkErrors(t, tt.name, i, err, tt.expectedError)
			if path != tt.expectedPath {
				t.Errorf("%s [%d] - incorrect returned path: received=%q | expected=%q",
					tt.name, i, path, tt.expectedPath)
			}
		})
	}

	if removeErr := os.RemoveAll(tempDir); removeErr != nil {
		t.Errorf("Error while deleting temporary dir: %q - %v", tempDir, removeErr)
	}
}

func TestBillingDownloader_Unzip(t *testing.T) {
	const (
		zipNoSuchFile = "./test/no-such-file.zip"
		zipDirPath    = "./test/test-dir.zip"
		zipFilePath   = "./test/test.txt.zip"
		zipFile       = "test.txt"
	)
	var tempDir = os.TempDir() + "/"

	t.Run("unzip directories and file", func(t *testing.T) {
		var expectedPath = tempDir
		var expectedError error

		d := &billingDownloader{}
		path, err := d.Unzip(zipDirPath, tempDir)
		if path != expectedPath {
			t.Errorf("Unzip invalid returned path: received=%q | expected=%q", path, expectedPath)
		}
		if err != expectedError {
			t.Errorf("Unzip invalid error returned: received=%+v | expected=%+v", err, expectedError)
		}
	})

	t.Run("unzip simple file", func(t *testing.T) {
		var expectedPath = tempDir + zipFile
		var expectedError error

		d := &billingDownloader{}
		path, err := d.Unzip(zipFilePath, tempDir)
		if path != expectedPath {
			t.Errorf("Unzip invalid returned path: received=%q | expected=%q", path, expectedPath)
		}
		if err != expectedError {
			t.Errorf("Unzip invalid error returned: received=%+v | expected=%+v", err, expectedError)
		}
	})

	t.Run("unzip no such file", func(t *testing.T) {
		var expectedPath = ""

		expectedError := errors.New("open ./test/no-such-file.zip: no such file or directory")
		d := &billingDownloader{}
		path, err := d.Unzip(zipNoSuchFile, tempDir)
		if path != expectedPath {
			t.Errorf("Unzip invalid returned path: received=%q | expected=%q", path, expectedPath)
		}
		if err != nil && expectedError.Error() != err.Error() {
			t.Errorf("Unzip invalid error returned: received=%+v | expected=%+v", err, expectedError)
		}
	})
}

func TestBillingDownloader_getAndCreateOutputPath(t *testing.T) {
	const (
		tempDownloadDir = "/billingDownloader/"
		tempFilename    = "test.csv.zip"
		givenFilename   = "bd.csv.zip"
	)
	var tempDir = os.TempDir() + tempDownloadDir
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		t.Errorf("Error while creating temporary dir: %q - %v", tempDir, err)
	}
	var tempFile = tempDir + tempFilename
	if _, err := os.Create(tempFile); err != nil {
		t.Errorf("Error while creating temporary file: %q - %v", tempFile, err)
	}

	tests := []struct {
		name          string
		destination   string
		expectedPath  string
		expectedError error
		expectExist   bool
	}{
		{name: "destination dir given exists",
			destination:   tempDir,
			expectedPath:  tempDir + givenFilename,
			expectedError: nil,
			expectExist:   false,
		},
		{name: "destination dir given does not exist",
			destination:   tempDir + "test/",
			expectedPath:  tempDir + "test/" + givenFilename,
			expectedError: nil,
			expectExist:   false,
		},
		{name: "destination file given exists",
			destination:   tempFile,
			expectedPath:  tempFile,
			expectedError: nil,
			expectExist:   true,
		},
		{name: "destination file given does not exist",
			destination:   tempDir + "otherTest.csv.zip",
			expectedPath:  tempDir + "otherTest.csv.zip",
			expectedError: nil,
			expectExist:   false,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &billingDownloader{}
			path, err := d.getAndCreateOutputPath(givenFilename, tt.destination)
			checkErrors(t, tt.name, i, err, tt.expectedError)
			if path != tt.expectedPath {
				t.Errorf("%s [%d] - incorrect path: received=%q | expected=%q",
					tt.name, i, path, tt.expectedPath)
			}
			_, statErr := os.Stat(path)
			if statErr != nil {
				if exists := os.IsNotExist(err); exists != tt.expectExist {
					t.Errorf("Error path should exist: %q", path)
				}
			} else if statErr == nil && tt.expectExist == false {
				t.Errorf("Error path shouldn't exist: %q", path)
			}
		})
	}

	if removeErr := os.RemoveAll(tempDir); removeErr != nil {
		t.Errorf("Error while deleting temporary dir: %q - %v", tempDir, removeErr)
	}
}
