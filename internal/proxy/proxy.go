package proxy

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"vbalancer/internal/peer"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

type Proxy struct {
	ctx              context.Context
	srv              net.Listener
	logger           *vlog.VLog
	peers            []*peer.Peer
	cfg              *Config
	currentPeerIndex *uint64
}

type Channel struct {
	from, to net.Conn
	ack      chan bool
}

func New(ctx context.Context, proxyPort string, cfg *Config, peers []*peer.Peer, logger *vlog.VLog) *Proxy {
	proxySrv, err := net.Listen("tcp", proxyPort)
	if err != nil {
		panic("connection error:" + err.Error())
	}

	var startPeerIndex uint64 = 0
	proxy := &Proxy{
		ctx:              ctx,
		srv:              proxySrv,
		logger:           logger,
		peers:            peers,
		cfg:              cfg,
		currentPeerIndex: &startPeerIndex,
	}

	return proxy
}

func (p *Proxy) Start(checkTimeAlive *peer.CheckTimeAlive) error {
	for _, peer := range p.peers {
		peer.CheckTimeAlive = checkTimeAlive
		peer.Logger = p.logger
		go peer.CheckIsAlive(p.ctx)
	}

	for {
		if conn, err := p.srv.Accept(); err == nil {
			conn.SetDeadline(time.Now().Add(time.Duration(p.cfg.TimeDeadLineMS) * time.Millisecond))
			go p.copyConn(conn)
		} else {
			p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(conn.RemoteAddr().String()), fmt.Sprintf("Accept failed, %v\n", err))
			fmt.Printf("Accept failed, %v\n", err)
		}
	}

}

func (p *Proxy) Shutdown(ctx context.Context) error {
	return nil
}

func (p *Proxy) copyConn(client net.Conn) {
	defer client.Close()
	peer, err := p.getNextPeer()

	if err != nil || peer == nil {
		// http.Error(w, "Proxy error", http.StatusServiceUnavailable)
		if peer != nil {
			p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(client.RemoteAddr().String()), vlog.ProxyHost(peer.Uri), err.Error())
		} else {
			p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(client.RemoteAddr().String()), err.Error())
		}
		return
	}

	dst, err := net.Dial("tcp", peer.Uri)
	if err != nil {
		p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(client.RemoteAddr().String()), fmt.Sprintf("Accept failed, %v\n", err))
		return
	}

	p.logger.Add(vlog.Debug,
		types.ResultOK,
		vlog.RemoteAddr(client.RemoteAddr().String()),
		vlog.ProxyHost(peer.Uri),
	)

	done := make(chan bool, 2)

	go func() {
		defer client.Close()
		defer dst.Close()
		io.Copy(dst, client)
		done <- true
	}()

	go func() {
		defer client.Close()
		defer dst.Close()
		io.Copy(client, dst)
		done <- true
	}()

	<-done
	<-done
}

func (p *Proxy) getNextPeer() (*peer.Peer, error) {
	var next int = 0
	if *p.currentPeerIndex >= uint64(len(p.peers)) {
		atomic.StoreUint64(p.currentPeerIndex, uint64(0))
	} else {
		next = p.nextIndex()
	}

	l := len(p.peers) + next
	for i := next; i < l; i++ {
		idx := i % len(p.peers)
		isAlive := p.peers[idx].IsAlive()
		if isAlive {
			atomic.StoreUint64(p.currentPeerIndex, uint64(idx))
			return p.peers[idx], nil
		}
	}
	return nil, fmt.Errorf("can't find active peers")
}

func (p *Proxy) nextIndex() int {
	return int(atomic.AddUint64(p.currentPeerIndex, uint64(1)) % uint64(len(p.peers)))
}

func (p *Proxy) copyHeader(source http.Header, dest *http.Header) {
	for n, v := range source {
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
}
