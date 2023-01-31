package proxy_test

import (
	"bytes"
	"context"
	"net"
	"testing"
	"vbalancer/internal/proxy"
	"vbalancer/mocks"
)

func TestCheckNewConnection(t *testing.T) {
	t.Parallel()

	proxySrv, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Error creating listener: %v", err)
	}
	defer proxySrv.Close()

	//nolint:exhaustivestruct,exhaustruct
	testProxy := &proxy.Proxy{
		Cfg:    &proxy.Config{DeadLineTimeMS: 100},
		Logger: &mocks.MockLogger{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go testProxy.CheckNewConnection(ctx, proxySrv)

	conn, err := net.Dial("tcp", proxySrv.Addr().String())
	if err != nil {
		t.Fatalf("Error dialing to proxy server: %v", err)
	}

	defer conn.Close()
}

func TestProxyDataCopy(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		clientData   []byte
		peerData     []byte
		expectedData []byte
	}{
		{
			clientData:   []byte("Hello World"),
			peerData:     []byte{},
			expectedData: []byte("Hello World"),
		},
	}

	for _, cases := range testCases {
		client := &mocks.MockConn{Data: cases.clientData, IsClient: true, Pos: 0} 
		peer := &mocks.MockConn{Data: cases.peerData, IsClient: false, Pos: 0}           

		//nolint:exhaustivestruct,exhaustruct
		proxyTest := &proxy.Proxy{Cfg: &proxy.Config{SizeCopyBufferIO: 64}} 

		proxyTest.ProxyDataCopy(client, peer)

		if !bytes.Equal(peer.Data, cases.expectedData) {
			t.Errorf("Expected %q, but got %q", cases.expectedData, client.Data)
		}
	}
}
