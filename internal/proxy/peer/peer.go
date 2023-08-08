package peer

import (
	"fmt"
	"net"
	"time"
)

// IPeer this is the interface that defines the methods for dialing a connection.
type IPeer interface {
	Dial(timeOut time.Duration, timeOutDeadLine time.Duration) (net.Conn, error)
	GetURI() string
}

// Peer this is the struct that implements the IPeer interface.
type Peer struct {
	Name   string `yaml:"name"`
	URI    string `yaml:"uri"`
}

// Dial dials a connection to a peer.
func (p *Peer) Dial(timeOut time.Duration, timeOutDeadLine time.Duration) (net.Conn, error) {
	connect, err := net.DialTimeout("tcp", p.GetURI(), timeOut)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	err = connect.SetDeadline(time.Now().Add(timeOutDeadLine))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return connect, nil
}

// GetURI return the URI of the peer.
func (p *Peer) GetURI() string {
	return p.URI
}
