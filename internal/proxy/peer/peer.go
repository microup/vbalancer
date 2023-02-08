package peer

import (
	"fmt"
	"net"
	"time"
	"vbalancer/internal/vlog"
)

type IPeer interface {
	Dial(timeOut time.Duration, timeOutDeadLine time.Duration) (net.Conn, error)
	GetURI() string
	SetLogger(vlog.ILog)
}

type Peer struct {
	Name   string `yaml:"name"`
	Proto  string `yaml:"proto"`
	URI    string `yaml:"uri"`
	logger vlog.ILog
}

func (p *Peer) Dial(timeOut time.Duration, timeOutDeadLine time.Duration) (net.Conn, error) {
	dst, err := net.DialTimeout("tcp", p.GetURI(), timeOut)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	err = dst.SetDeadline(time.Now().Add(timeOutDeadLine))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return dst, nil
}

func (p *Peer) SetLogger(value vlog.ILog) {
	p.logger = value
}

func (p *Peer) GetURI() string {
	return p.URI
}
