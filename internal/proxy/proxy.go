package proxy

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"vbalancer/internal/peer"
	"vbalancer/internal/vlog"
)

type Proxy struct {
	ctx                     context.Context
	proxyServer             *http.Server
	logger                  *vlog.VLog
	peers                   []*peer.Peer
	cfg                     *Config
	currentPeerIndex        *uint64
}

func New(ctx context.Context, proxyPort string, cfg *Config, peers []*peer.Peer, logger *vlog.VLog) *Proxy {

	httpServer := &http.Server{
		Addr:              proxyPort,
		TLSConfig:         nil,
		ReadTimeout:       time.Duration(cfg.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(cfg.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(cfg.IdleTimeout) * time.Second,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	var startPeerIndex uint64 = 0
	proxy := &Proxy{
		ctx:              ctx,
		proxyServer:      httpServer,
		logger:           logger,
		peers:            peers,
		cfg:              cfg,
		currentPeerIndex: &startPeerIndex,
	}

	httpServer.Handler = http.HandlerFunc(proxy.ProxyHandler)

	return proxy
}

func (p *Proxy) Start(checkTimeAlive *peer.CheckTimeAlive) error {
	for _, peer := range p.peers {
		peer.CheckTimeAlive = checkTimeAlive
		peer.Logger = p.logger
		go peer.CheckIsAlive(p.ctx)
	}

	return p.proxyServer.ListenAndServe()
}

func (p *Proxy) Shutdown(ctx context.Context) error {
	return p.proxyServer.Shutdown(ctx)
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
