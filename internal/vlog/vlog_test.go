package vlog_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

//nolint:paralleltest
func TestVlogAdd(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")

		return
	}

	helperVlogAdd(t)
}

func helperVlogAdd(t *testing.T) {
	t.Helper()

	cfg := &vlog.Config{
		DirLog:         "./logs/",
		FileSize:       1000,
		APIShowRecords: 5,
	}

	vLog, err := vlog.New(cfg)
	if err != nil {
		t.Fatalf("unexpected error creating VLog: %v", err)
	}

	vLog.Add(types.Debug, "test msg")
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
