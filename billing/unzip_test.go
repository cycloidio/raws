package billing

import (
	"testing"
	"os"
	"errors"
)

func TestUnzip(t *testing.T) {
	const (
		zipNoSuchFile = "./test/no-such-file.zip"
		zipDirPath = "./test/test-dir.zip"
		zipFilePath = "./test/test.txt.zip"
		zipFile = "test.txt"
	)
	var tempDir string = os.TempDir() + "/"

	t.Run("unzip directories and file", func(t *testing.T) {
		var expectedPath string = tempDir
		var expectedError error = nil

		path, err := Unzip(zipDirPath, tempDir)
		if path != expectedPath {
			t.Errorf("Unzip invalid returned path: received=%q | expected=%q", path, expectedPath)
		}
		if err != expectedError {
			t.Errorf("Unzip invalid error returned: received=%+v | expected=%+v", err, expectedError)
		}
	})

	t.Run("unzip simple file", func(t *testing.T) {
		var expectedPath string = tempDir + zipFile
		var expectedError error = nil

		path, err := Unzip(zipFilePath, tempDir)
		if path != expectedPath {
			t.Errorf("Unzip invalid returned path: received=%q | expected=%q", path, expectedPath)
		}
		if err != expectedError {
			t.Errorf("Unzip invalid error returned: received=%+v | expected=%+v", err, expectedError)
		}
	})

	t.Run("unzip no such file", func(t *testing.T) {
		var expectedPath string = ""
		var expectedError error = errors.New("open ./test/no-such-file.zip: no such file or directory")

		path, err := Unzip(zipNoSuchFile, tempDir)
		if path != expectedPath {
			t.Errorf("Unzip invalid returned path: received=%q | expected=%q", path, expectedPath)
		}
		if err != nil && expectedError.Error() != err.Error() {
			t.Errorf("Unzip invalid error returned: received=%+v | expected=%+v", err, expectedError)
		}
	})
}
