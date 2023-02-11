package proxy_test

import (
	"context"
	"net"
	"testing"

	"vbalancer/internal/proxy"
	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/proxy/peers"
	"vbalancer/mocks"
)

func TestCheckNewConnection(t *testing.T) {
	t.Parallel()

	proxySrv, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Error creating listener: %v", err)
	}
	defer proxySrv.Close()

	logger := &mocks.MockLogger{}

	listPeer := make([]peer.IPeer, 0)
	testPeer := peer.Peer{
		Name:  "test peer",
		URI:   "127.0.0.1:0",
	}
	listPeer = append(listPeer, &testPeer)

	//nolint:exhaustivestruct,exhaustruct
	testProxy := &proxy.Proxy{
		Cfg: &proxy.Config{},

		Logger: logger,
		Peers:  peers.New(listPeer),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go testProxy.AcceptConnections(ctx, proxySrv)

	conn, err := net.Dial("tcp", proxySrv.Addr().String())
	if err != nil {
		t.Fatalf("Error dialing to proxy server: %v", err)
	}

	defer conn.Close()
}
