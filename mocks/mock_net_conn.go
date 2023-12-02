package mocks

import (
	"io"
	"net"
	"sync"
	"time"
)

type MockConn struct {
	Data     []byte
	Pos      int
	IsClient bool
	mu       sync.Mutex // guards
}

func (c *MockConn) Read(inBytes []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Pos >= len(c.Data) {
		return 0, io.EOF
	}

	n := copy(inBytes, c.Data[c.Pos:])
	c.Pos += n

	return n, nil
}

func (c *MockConn) Write(inBytes []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Data = append(c.Data, inBytes...)
	c.Pos += len(c.Data)

	return len(c.Data), nil
}

func (c *MockConn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return nil
}

func (c *MockConn) LocalAddr() net.Addr {
	c.mu.Lock()
	defer c.mu.Unlock()

	return nil
}

func (c *MockConn) RemoteAddr() net.Addr {
	c.mu.Lock()
	defer c.mu.Unlock()

	return nil
}

func (c *MockConn) SetDeadline(time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return nil
}

func (c *MockConn) SetReadDeadline(time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return nil
}

func (c *MockConn) SetWriteDeadline(time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return nil
}
