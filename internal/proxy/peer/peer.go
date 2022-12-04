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

type Peer struct {
	Name           string `yaml:"name"`
	Proto          string `yaml:"proto"`
	URI            string `yaml:"uri"`
	CheckTimeAlive *CheckTimeAlive
	Alive          bool
	Logger         *vlog.VLog
	urlPeer        *url.URL
	mux            sync.RWMutex
}

func (p *Peer) setAlive(alive bool) {
	p.mux.Lock()
	p.Alive = alive
	p.mux.Unlock()
}

func (p *Peer) IsAlive() bool {
	p.mux.RLock()
	alive := p.Alive
	p.mux.RUnlock()

	return alive
}

func (p *Peer) CheckIsAlive(ctx context.Context) {
	if p.urlPeer == nil {
		p.urlPeer, _ = url.Parse(fmt.Sprintf("%s://%s", p.Proto, p.URI))
	}

	timeout := time.Duration(p.CheckTimeAlive.TimeCheck) * time.Second

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
		time.Sleep(time.Duration(p.CheckTimeAlive.WaitTimeCheck) * time.Second)
	}
}
