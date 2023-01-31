package proxy_test

import (
	"context"
	"net"
	"testing"
	"vbalancer/internal/proxy"
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
		Cfg: &proxy.Config{DeadLineTimeMS: 100},
		Logger: &MockLogger{
			t: t,
		},
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

type MockLogger struct {
	t *testing.T
}

func (m *MockLogger) Add(values ...interface{}) {

}

func (m *MockLogger) Close() error {
	return nil
}