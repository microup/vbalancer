package proxy_test

import (
	"bytes"
	"context"
	"io"
	"net"
	"testing"
	"time"
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

func TestProxyDataCopy(t *testing.T) {
	t.Parallel()

	// Create a test case with sample input and expected output data
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
		{
			clientData:   []byte("123456789"),
			peerData:     []byte{},
			expectedData: []byte("123456789"),
		},
	}

	for _, cases := range testCases {
		client := &FakeConn{data: cases.clientData} //nolint:exhaustivestruct,exhaustruct
		peer := &FakeConn{data: []byte{}} //nolint:exhaustivestruct,exhaustruct

		proxyTest := &proxy.Proxy{Cfg: &proxy.Config{SizeCopyBufferIO: 1024}} //nolint:exhaustivestruct,exhaustruct

		proxyTest.ProxyDataCopy(client, peer)

		done := make(chan struct{}, 2)
		proxyTest.CopyDataClientToPeer(peer, client, done)
		done <- struct{}{}

		if !bytes.Equal(peer.data, cases.expectedData) {
			t.Errorf("Expected %q, but got %q", cases.expectedData, client.data)
		}
	}
}

type MockLogger struct {
	t *testing.T
}

func (m *MockLogger) Add(values ...interface{}) {

}

func (m *MockLogger) Close() error {
	return nil
}


// fakeConn is a fake implementation of the net.Conn interface
type FakeConn struct {
	data []byte
	pos  int
}

func (c *FakeConn) Read(b []byte) (int, error) {
	if c.pos >= len(c.data) {
		return 0, io.EOF
	}

	n := copy(b, c.data[c.pos:])
	c.pos += n

	return n, nil
}

func (c *FakeConn) Write(b []byte) (int, error) {
	c.data = append(c.data, b...)
	c.pos += len(b)

	return len(b), nil
}

func (c *FakeConn) Close() error {
	return nil
}

func (c *FakeConn) LocalAddr() net.Addr {
	return nil
}

func (c *FakeConn) RemoteAddr() net.Addr {
	return nil
}

func (c *FakeConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *FakeConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *FakeConn) SetWriteDeadline(t time.Time) error {
	return nil
}