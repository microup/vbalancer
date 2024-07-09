package vlog_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"vbalancer/internal/config"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestVlogAdd is a test case for vlog.Add.
func TestVlogAdd(t *testing.T) {
	t.Parallel()

	helperVlogAdd(t)
}

func helperVlogAdd(t *testing.T) {
	t.Helper()

	if os.Getenv("CI") != "" {
		return
	}

	cfg := &config.Log{
		DirLog:         "./logs/",
		FileSizeMB:     1,
		APIShowRecords: 5,
	}

	vLog := vlog.New(cfg)

	err := vLog.Init()

	require.NoError(t, err)

	vLog.Add(types.Debug, "test msg")
	time.Sleep(1 * time.Second)

	assert.Equal(t, 1, vLog.GetCountRecords(), "expected count of log records to be 1, got %d", vLog.GetCountRecords())

	err = vLog.Close()

	require.NoError(t, err)

	absolutePath, _ := filepath.Abs(cfg.DirLog)
	err = os.RemoveAll(absolutePath)

	require.NoErrorf(t, err, "unexpected error delete tempore dir")
}
