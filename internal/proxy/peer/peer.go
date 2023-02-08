package peer

import (
	"vbalancer/internal/vlog"
)

type IPeer interface {
	GetURI() string
	SetLogger(vlog.ILog)
}

type Peer struct {
	Name   string `yaml:"name"`
	Proto  string `yaml:"proto"`
	URI    string `yaml:"uri"`
	logger vlog.ILog
}

func (p *Peer) SetLogger(value vlog.ILog) {
	p.logger = value
}

func (p *Peer) GetURI() string {
	return p.URI
}
