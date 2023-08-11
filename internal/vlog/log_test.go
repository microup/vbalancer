package vlog_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"vbalancer/internal/config"
	"vbalancer/internal/vlog"
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
	if err != nil {
		t.Fatalf("unexpected error creating VLog: %v", err)
	}

	vLog.Add(vlog.Debug, "test msg")
	time.Sleep(1 * time.Second)

	if vLog.GetCountRecords() != 1 {
		t.Errorf("expected count of log records to be 1, got %d", vLog.GetCountRecords())
	}

	err = vLog.Close()
	if err != nil {
		t.Fatalf("unexpected close log file: %v", err)
	}

	absolutePath, _ := filepath.Abs(cfg.DirLog)
	err = os.RemoveAll(absolutePath)

	if err != nil {
		t.Fatalf("unexpected error delete tempore dir: %v", err)
	}
}
