package peer

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"
	"vbalancer/internal/vlog"
)

type IPeer interface {
	IsAvailable() bool
	CheckAvailability(context.Context)
	GetURI() string
	SetAvailabilityCheckInterval(*CheckTimeAlive)
	SetLogger(vlog.ILog)
}

type Peer struct {
	Name           string `yaml:"name"`
	Proto          string `yaml:"proto"`
	URI            string `yaml:"uri"`
	alive          bool
	checkTimeAlive *CheckTimeAlive
	logger         vlog.ILog
	urlPeer        *url.URL
	Mu             *sync.RWMutex
}

func (p *Peer) CheckAvailability(ctx context.Context) {
	if p.urlPeer == nil {
		p.urlPeer, _ = url.Parse(fmt.Sprintf("%s://%s", p.Proto, p.URI))
	}

	timeout := time.Duration(p.checkTimeAlive.TimeCheck) * time.Second

	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := net.DialTimeout("tcp", p.urlPeer.Host, timeout)

			if err != nil {
				p.setAlive(false)
			} else {
				_ = conn.Close()
				p.setAlive(true)
			}
		}
		time.Sleep(time.Duration(p.checkTimeAlive.WaitTimeCheck) * time.Second)
	}
}

func (p *Peer) SetAvailabilityCheckInterval(value *CheckTimeAlive) {
	p.checkTimeAlive = value
}

func (p *Peer) SetLogger(value vlog.ILog) {
	p.logger = value
}

func (p *Peer) IsAvailable() bool {
	p.Mu.RLock()
	alive := p.alive
	p.Mu.RUnlock()

	return alive
}

func (p *Peer) GetURI() string {
	return p.URI
}

func (p *Peer) setAlive(value bool) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	p.alive = value
}
