//nolint:testpackage // the need to test the private method updatePort() in the proxy struct.
package proxy

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/proxy/peers"
	"vbalancer/internal/types"
	"vbalancer/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProxyServer - this is the `TestProxyServer` function of the proxy server.
func TestCheckNewConnection(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port := "18880"
	uri := fmt.Sprintf("127.0.0.1:%s", port)

	listener, err := net.Listen("tcp", uri)

	require.NoErrorf(t, err, "error creating listener")

	defer listener.Close()

	logger := &mocks.MockLogger{}

	listPeer := []peer.Peer{
		{
			Name: "test client",
			URI:  "127.0.0.1:18881",
		},
	}

	testProxy := &Proxy{
		Logger:                logger,
		Port:                  port,
		ClientDeadLineTime:    1,
		PeerConnectionTimeout: 1,
		MaxCountConnection:    100,
		Peers:                 nil,
		Rules:                 nil,
	}

	resultCode := testProxy.updatePort()
	assert.Equal(t, types.ResultOK, resultCode, "name: `%s`")

	testProxy.Peers = &peers.Peers{
		TimeToEvictNotResponsePeers: 0,
		Peers:                       nil,
	}

	err = testProxy.Peers.Init(ctx, listPeer)

	require.NoErrorf(t, err, "can't init peers")

	go testProxy.AcceptConnections(ctx, listener)

	time.Sleep(2 * time.Second)

	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s", uri), nil)
	if err != nil {
		require.NoErrorf(t, err, "http request")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		require.NoErrorf(t, err, "http client do")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		assert.Equal(t, types.ResultOK, resp.StatusCode, "name: `%s`")
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		require.NoErrorf(t, err, "io readall")
	}
}

// TestGetProxyPort tests the UpdatePort function.
// It validates UpdatePort handles invalid environment variable values,
// default values, and valid custom environment variable values correctly.
func TestGetProxyPort(t *testing.T) {
	testCases := []struct {
		port      string
		name      string
		envVar    string
		want      types.ResultCode
		wantValue string
	}{
		{
			name:      "set port `1234`",
			port:      "1234",
			envVar:    ":",
			want:      types.ResultOK,
			wantValue: ":1234",
		},
		{
			name:      "empty env var, got DefaultPort",
			port:      "",
			envVar:    "",
			want:      types.ResultOK,
			wantValue: fmt.Sprintf(":%s", types.DefaultProxyPort),
		},
		{
			name:      "valid proxy port from env var",
			port:      "",
			envVar:    "8080",
			want:      types.ResultOK,
			wantValue: ":8080",
		},
		{
			name:      "empty proxy port from default value",
			envVar:    " ",
			want:      types.ErrEmptyValue,
			wantValue: ":",
			port:      ":",
		},
		{
			name:      "empty proxy port from default value",
			envVar:    "          ",
			want:      types.ErrEmptyValue,
			wantValue: ":",
			port:      ":",
		},
	}

	prx := &Proxy{
		Logger:                nil,
		Port:                  "",
		ClientDeadLineTime:    10,
		PeerConnectionTimeout: 10,
		MaxCountConnection:    100,
		Peers:                 nil,
		Rules:                 nil,
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			prx.Port = testCase.port

			t.Setenv(types.ProxyPort, testCase.envVar)

			result := prx.updatePort()

			assert.Equal(t, testCase.want, result, "name: `%s`")

			assert.Equal(t, testCase.wantValue, prx.Port, "name: `%s`")
		})
	}
}
