package proxy_test

import (
	"context"
	"net"
	"testing"

	"vbalancer/internal/config"
	"vbalancer/internal/proxy"
	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/proxy/peers"
	"vbalancer/mocks"
)

// TestProxyServer - this is the `TestProxyServer` function of the proxy server.
func TestCheckNewConnection(t *testing.T) {
	t.Parallel()

	proxySrv, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Error creating listener: %v", err)
	}
	defer proxySrv.Close()

	logger := &mocks.MockLogger{}

	listPeer := make([]peer.Peer, 0)
	testPeer := peer.Peer{
		Name: "test peer",
		URI:  "127.0.0.1:8080",
	}
	listPeer = append(listPeer, testPeer)

	//nolint:exhaustivestruct,exhaustruct
	testProxy := &proxy.Proxy{
		Logger: logger,
		Cfg: &config.Proxy{
			DefaultPort: "8080",
			ClientDeadLineTime: 10,
			PeerHostTimeOut: 10,
			PeerHostDeadLine: 10,
			MaxCountConnection: 100,
			CountDialAttemptsToPeer: 10,
		},
	}

	testProxy.Peers = peers.New()
	
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
