package peer

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/microup/vbalancer/internal/types"
	"github.com/microup/vbalancer/internal/vlog"
)

type Peer struct {
	Name           string `yaml:"Name"`
	Proto          string `yaml:"Proto"`
	Uri            string `yaml:"URI"`
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

func (p *Peer) IsAlive() (alive bool) {
	p.mux.RLock()
	alive = p.Alive
	p.mux.RUnlock()
	return
}

func (p *Peer) CheckIsAlive(ctx context.Context) {
	if p.urlPeer == nil {
		p.urlPeer, _ = url.Parse(fmt.Sprintf("%s://%s", p.Proto, p.Uri))
	}

	timeout := time.Duration(p.CheckTimeAlive.TimeCheck) * time.Second

	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := net.DialTimeout("tcp", p.urlPeer.Host, timeout)

			if err != nil {
				p.Logger.Add(vlog.Debug,
					types.ResultCode(types.StatusInternalServerError),
					vlog.RemoteAddr(p.urlPeer.Host),
					vlog.ProxyHost(p.urlPeer.Host),
					vlog.ProxyProto(p.Proto),
					err.Error())
				p.setAlive(false)
			} else {
				_ = conn.Close()
				p.setAlive(true)
			}
		}
		time.Sleep(time.Duration(p.CheckTimeAlive.WaitTimeCheck) * time.Second)		
	}
}
