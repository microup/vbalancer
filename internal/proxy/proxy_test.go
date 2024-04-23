//nolint:testpackage
package proxy

import (
	"context"
	"fmt"
	"net"
	"testing"

	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/proxy/peers"
	"vbalancer/internal/types"
	"vbalancer/mocks"

	"github.com/stretchr/testify/assert"
)

// TestProxyServer - this is the `TestProxyServer` function of the proxy server.
func TestCheckNewConnection(t *testing.T) {
	t.Parallel()

	proxySrv, err := net.Listen("tcp", "127.0.0.1:18880")

	assert.Nil(t, err, "error creating listener")

	defer proxySrv.Close()

	logger := &mocks.MockLogger{}

	listPeer := make([]peer.Peer, 0)
	testPeer := peer.Peer{
		Name: "test peer",
		URI:  "127.0.0.1:18880",
	}
	listPeer = append(listPeer, testPeer)

	//nolint:exhaustivestruct,exhaustruct
	testProxy := &Proxy{
		Logger:                logger,
		Port:                  "18880",
		ClientDeadLineTime:    10,
		PeerConnectionTimeout: 10,
		PeerHostDeadLine:      10,
		MaxCountConnection:    100,
	}

	resultCode := testProxy.updatePort()
	assert.Equal(t, resultCode, types.ResultOK, "name: `%s`")

	//nolint:exhaustivestruct,exhaustruct
	testProxy.Peers = &peers.Peers{}

	err = testProxy.Peers.Init(context.Background(), listPeer)

	assert.Nil(t, err, "can't init peers")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go testProxy.AcceptConnections(ctx, proxySrv)

	conn, err := net.Dial("tcp", proxySrv.Addr().String())

	assert.Nil(t, err, "dialing to proxy server")

	defer conn.Close()
}

// TestGetProxyPort tests the UpdatePort function.
// It validates UpdatePort handles invalid environment variable values,
// default values, and valid custom environment variable values correctly.
//
//nolint:paralleltest
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
	//nolint:exhaustivestruct,exhaustruct
	prx := &Proxy{}

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
