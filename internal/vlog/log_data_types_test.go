package vlog_test

import (
	"strconv"
	"strings"
	"testing"
	"time"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

// TestBuildRecord.
func TestBuildRecord(t *testing.T) {
	t.Parallel()

	typeLog := vlog.Info
	resultCode := types.ResultCode(200)
	remoteAddr := vlog.RemoteAddr("192.168.1.40")
	peerAddr := vlog.PeerAddr("127.0.0.1:8081")
	valuesStr := "value1;value2"

	// Call the function with the inputs
	actualTypeLog, actualRecord := vlog.BuildRecord(typeLog, resultCode, remoteAddr, peerAddr,valuesStr)

	// Define the expected output
	expectedTypeLog := vlog.Info
	expectedRecord := "INFO;200;192.168.1.40;127.0.0.1:8081;value1;value2;"

	// Assert that the actual output is as expected
	if actualTypeLog != expectedTypeLog {
		t.Errorf("Expected typeLog %d but got %d", expectedTypeLog, actualTypeLog)
	}

	parts := strings.Split(actualRecord, ";")
	dateTime := parts[:2]

	_, err := time.Parse("2006-01-02;15:04:05", strings.Join(dateTime, ";"))
	if err != nil {
		t.Errorf("Expected valid date and time, but got %v", dateTime)
	}

	// Check that the fourth element of parts is equal to resultCode
	if parts[3] != strconv.Itoa(int(resultCode.ToUint())) {
		t.Errorf("Expected resultCode %d, but got %v", resultCode, parts[3])
	}

	resultStr := strings.Join(parts[2:], ";")

	if expectedRecord != resultStr {
		t.Errorf("Expected record %s but got %s", expectedRecord, resultStr)
	}
}
