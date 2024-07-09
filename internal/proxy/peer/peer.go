package peer

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Peerer this is the interface that defines the methods for dialing a connection.
type Peerer interface {
	GetURI() string
	Dial(ctx context.Context, timeOut time.Duration, timeOutDeadLine time.Duration) (net.Conn, error)
}

// Peer this is the struct that implements the IPeer interface.
type Peer struct {
	Name string `yaml:"name"`
	URI  string `yaml:"uri"`
}

// Dial dials a connection to a peer.
func (p *Peer) Dial(timeOut time.Duration) (net.Conn, error) {
	connect, err := net.DialTimeout("tcp", p.GetURI(), timeOut)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return connect, nil
}

// GetURI return the URI of the peer.
func (p *Peer) GetURI() string {
	return p.URI
}
