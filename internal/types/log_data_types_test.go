package types_test

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"vbalancer/internal/types"
	"vbalancer/internal/vlog"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuildRecord.
func TestBuildRecord(t *testing.T) {
	t.Parallel()

	typeLog := types.Info
	resultCode := types.ResultCode(200)
	remoteAddr := types.RemoteAddr("192.168.1.40")
	peerAddr := types.PeerAddr("127.0.0.1:8081")
	valuesStr := "value1;value2"

	// Call the function with the inputs
	actualTypeLog, actualRecord := vlog.BuildRecord(typeLog, resultCode, remoteAddr, peerAddr, valuesStr)

	// Define the expected output
	expectedTypeLog := types.Info
	expectedRecord := "INFO;200;192.168.1.40;127.0.0.1:8081;value1;value2;"

	assert.Equal(t, expectedTypeLog, actualTypeLog)

	parts := strings.Split(actualRecord, ";")
	dateTime := parts[:2]

	_, err := time.Parse("2006-01-02;15:04:05", strings.Join(dateTime, ";"))

	require.NoErrorf(t, err, "expected valid date and time, but got %v", dateTime)

	// Check that the fourth element of parts is equal to resultCode
	assert.Equal(t, strconv.Itoa(int(resultCode.ToUint())), parts[3])

	resultStr := strings.Join(parts[2:], ";")

	assert.Equal(t, expectedRecord, resultStr)
}
