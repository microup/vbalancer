package core_test

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"vbalancer/internal/core"
)

// TestArchive this is a test function for `archive()`.
//nolint:paralleltest
func TestArchiveFile(t *testing.T) {
	helperArchiveFile(t)
}

//nolint:funlen,cyclop
func helperArchiveFile(t *testing.T) {
	t.Helper()

	if os.Getenv("CI") != "" {
		return
	}

	fileName := "test_file.csv"
	extension := ".zip"

	testFile, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	defer os.Remove(fileName)

	_, err = testFile.Write([]byte("test data"))
	if err != nil {
		t.Fatalf("Failed to write data to test file: %v", err)
	}

	err = testFile.Close()
	if err != nil {
		t.Fatalf("Failed to close test file: %v", err)
	}

	err = core.ArchiveFile(fileName, extension)
	if err != nil {
		t.Fatalf("Archiving failed: %v", err)
	}

	archivedFile := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + extension
	if _, err = os.Stat(archivedFile); os.IsNotExist(err) {
		t.Fatalf("Archived file does not exist: %v", err)
	}

	defer os.Remove(archivedFile)

	zipFile, err := zip.OpenReader(archivedFile)
	if err != nil {
		t.Fatalf("Failed to open archived file: %v", err)
	}
	defer zipFile.Close()

	if len(zipFile.File) != 1 {
		t.Fatalf("Unexpected number of files in archived file: %d", len(zipFile.File))
	}

	fileInArchive := zipFile.File[0]
	zipFileContent, err := fileInArchive.Open()

	if err != nil {
		t.Fatalf("Failed to open file in archive: %v", err)
	}
	defer zipFileContent.Close()

	data, err := io.ReadAll(zipFileContent)
	if err != nil {
		t.Fatalf("Failed to read data from file in archive: %v", err)
	}

	if string(data) != "test data" {
		t.Fatalf("Unexpected data in archived file: %s", string(data))
	}
}
