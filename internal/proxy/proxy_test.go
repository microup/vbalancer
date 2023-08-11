package proxy_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"

	"vbalancer/internal/proxy"
	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/proxy/peers"
	"vbalancer/internal/types"
	"vbalancer/mocks"
)

// TestProxyServer - this is the `TestProxyServer` function of the proxy server.
func TestCheckNewConnection(t *testing.T) {
	t.Parallel()

	proxySrv, err := net.Listen("tcp", "127.0.0.1:18880")
	if err != nil {
		t.Fatalf("Error creating listener: %v", err)
	}
	defer proxySrv.Close()

	logger := &mocks.MockLogger{}

	listPeer := make([]peer.Peer, 0)
	testPeer := peer.Peer{
		Name: "test peer",
		URI:  "127.0.0.1:18880",
	}
	listPeer = append(listPeer, testPeer)

	//nolint:exhaustivestruct,exhaustruct
	testProxy := &proxy.Proxy{
		Logger:                     logger,
		Port:                       "18880",
		ClientDeadLineTime:         10,
		PeerHostTimeOut:            10,
		PeerHostDeadLine:           10,
		MaxCountConnection:         100,
		CountMaxDialAttemptsToPeer: 10,
	}

	resultCode := testProxy.UpdatePort()
	if resultCode != types.ResultOK {
		t.Fatalf("can't update proxy port: %d", resultCode)
	}

	//nolint:exhaustivestruct,exhaustruct
	testProxy.Peers = &peers.Peers{}

	err = testProxy.Peers.Init(listPeer)
	if err != nil {
		t.Fatalf("can't init peers: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go testProxy.AcceptConnections(ctx, proxySrv)

	conn, err := net.Dial("tcp", proxySrv.Addr().String())
	if err != nil {
		t.Fatalf("dialing to proxy server: %v", err)
	}

	defer conn.Close()
}

// TestGetProxyPort tests the UpdatePort function.
// It validates UpdatePort handles invalid environment variable values,
// default values, and valid custom environment variable values correctly.
//
//nolint:funlen
func TestGetProxyPort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		port       string
		name       string
		envVar     string
		expected   types.ResultCode
		checkValue string
	}{
		{
			name:       "set port `1234`",
			port:       "1234",
			envVar:     ":",
			expected:   types.ResultOK,
			checkValue: ":1234",
		},
		{
			name:       "empty env var, got DefaultPort",
			envVar:     "",
			expected:   types.ResultOK,
			checkValue: fmt.Sprintf(":%s", types.DefaultProxyPort),
		},
		{
			name:       "valid proxy port from env var",
			envVar:     "8080",
			expected:   types.ResultOK,
			checkValue: ":8080",
		},
		{
			name:       "empty proxy port from default value",
			envVar:     " ",
			expected:   types.ErrEmptyValue,
			checkValue: ":",
			port:       ":",
		},
		{
			name:       "empty proxy port from default value",
			envVar:     "          ",
			expected:   types.ErrEmptyValue,
			checkValue: ":",
			port:       ":",
		},
	}

	//nolint:exhaustivestruct,exhaustruct
	prx := &proxy.Proxy{}

	for _, test := range tests {
		prx.Port = test.port

		os.Clearenv()
		os.Setenv("ProxyPort", test.envVar)

		result := prx.UpdatePort()
		if result != test.expected {
			t.Fatalf("name: %s, expected result %v, got %v", test.name, test.expected, result)
		}

		if prx.Port != test.checkValue {
			t.Fatalf("name: %s, expected value %s, got %s", test.name, test.checkValue, prx.Port)
		}
	}
}
