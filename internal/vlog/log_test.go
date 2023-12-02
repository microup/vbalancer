package vlog_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"vbalancer/internal/config"
	"vbalancer/internal/vlog"

	"github.com/stretchr/testify/assert"
)

// TestVlogAdd is a test case for vlog.Add.
//
//nolint:paralleltest
func TestVlogAdd(t *testing.T) {
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

	assert.Nil(t, err)

	vLog.Add(vlog.Debug, "test msg")
	time.Sleep(1 * time.Second)

	assert.Equal(t, vLog.GetCountRecords(), 1, "expected count of log records to be 1, got %d", vLog.GetCountRecords())

	err = vLog.Close()

	assert.Nil(t, err)

	absolutePath, _ := filepath.Abs(cfg.DirLog)
	err = os.RemoveAll(absolutePath)

	assert.Nil(t, err, "unexpected error delete tempore dir")
}
