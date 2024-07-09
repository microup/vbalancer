package core_test

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"vbalancer/internal/core"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestArchive this is a test function for `archive()`.
//
//nolint:paralleltest // this why need to run once
func TestArchiveFile(t *testing.T) {
	helperArchiveFile(t)
}

func helperArchiveFile(t *testing.T) {
	t.Helper()

	if os.Getenv("CI") != "" {
		return
	}

	fileName := "test_file.csv"
	extension := ".zip"

	testFile, err := os.Create(fileName)

	require.NoErrorf(t, err, "failed to create test file")

	defer os.Remove(fileName)

	_, err = testFile.WriteString("test data")

	require.NoErrorf(t, err, "failed to write data to test file")

	err = testFile.Close()

	require.NoErrorf(t, err, "failed to close test file")

	err = core.ArchiveFile(fileName, extension)

	require.NoErrorf(t, err, "archiving failed")

	archivedFile := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + extension
	if _, err = os.Stat(archivedFile); os.IsNotExist(err) {
		assert.FailNow(t, "archived file does not exist", err)
	}

	require.NoErrorf(t, err, "archiving failed")

	defer os.Remove(archivedFile)

	zipFile, err := zip.OpenReader(archivedFile)
	if err != nil {
		assert.FailNow(t, "failed to open archived file: %v", err)
	}
	defer zipFile.Close()

	fileInArchive := zipFile.File[0]
	zipFileContent, err := fileInArchive.Open()

	require.NoErrorf(t, err, "failed to open file in archive")

	defer zipFileContent.Close()

	data, err := io.ReadAll(zipFileContent)

	require.NoErrorf(t, err, "failed to read data from file in archive")

	assert.Equal(t, "test data", string(data), "unexpected data in archived file")
}
