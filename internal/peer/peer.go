package peer

import (
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/microup/vbalancer/internal/types"
	"github.com/microup/vbalancer/internal/vlog"
)

type Peer struct {
	Name    string `yaml:"Name"`
	Proto   string `yaml:"Proto"`
	Uri     string `yaml:"URI"`
	Logger  *vlog.VLog
	urlPeer *url.URL
	mux     sync.RWMutex
}

func (p *Peer) IsBackendAlive() bool {
	if p.urlPeer == nil {
		p.urlPeer, _ = url.Parse(fmt.Sprintf("%s://%s", p.Proto, p.Uri))
	}
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", p.urlPeer.Host, timeout)
	if err != nil {
		p.Logger.Add(vlog.Debug, types.ResultCode(types.StatusInternalServerError), 
			 vlog.RemoteAddr(p.urlPeer.Host), 
			 vlog.ProxyHost(p.urlPeer.Host), 
			 vlog.ProxyProto(p.Proto), 
			 err.Error())
		return false
	}
	_ = conn.Close()
	return true
}
