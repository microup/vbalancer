package types_test

import (
	"strconv"
	"strings"
	"testing"
	"time"
	"vbalancer/internal/types"
)

func TestBuildRecord(t *testing.T) {
	t.Parallel()

	typeLog := types.Info
	resultCode := types.ResultCode(200)
	remoteAddr := types.RemoteAddr("127.0.0.1")
	clientHost := types.ClientHost("client.com")
	clientMethod := types.ClientMethod("GET")
	clientProto := types.ClientProto("HTTP/1.1")
	clientURI := types.ClientURI("/client")
	proxyHost := types.ProxyHost("proxy.com")
	proxyMethod := types.ProxyMethod("GET")
	proxyProto := types.ProxyProto("HTTP/1.1")
	proxyURI := types.ProxyURI("/proxy")
	valuesStr := "value1;value2"

	// Call the function with the inputs
	actualTypeLog, actualRecord := types.BuildRecord(
		typeLog, resultCode, remoteAddr, clientHost, clientMethod,
		clientProto, clientURI, proxyHost, proxyMethod, proxyProto, proxyURI,
		valuesStr)

	// Define the expected output
	expectedTypeLog := types.Info
	expectedRecord := "INFO;200;127.0.0.1;client.com;GET;HTTP/1.1;/client;GET;HTTP/1.1;" +
		"proxy.com;/proxy;value1;value2"

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
