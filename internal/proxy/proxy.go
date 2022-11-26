package proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/microup/vbalancer/internal/peer"
	"github.com/microup/vbalancer/internal/types"
	"github.com/microup/vbalancer/internal/vlog"
)

type Proxy struct {
	ctx              context.Context
	proxyServer      *http.Server
	logger           *vlog.VLog
	peers            []*peer.Peer
	cfg              *Config
	currentPeerIndex *uint64
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

func (p *Proxy) ProxyHandler(w http.ResponseWriter, r *http.Request) {
	p.logger.Add(vlog.Debug, types.ResultOK, vlog.RemoteAddr(r.RemoteAddr), vlog.ClientHost(r.Host), vlog.ClientMethod(r.Method), vlog.ClientProto(r.Proto), vlog.ClientURI(r.RequestURI))

	if len(p.peers) == 0 {
		http.Error(w, "Proxy error", http.StatusInternalServerError)
		p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(r.RemoteAddr), vlog.ClientHost(r.Host), vlog.ClientMethod(r.Method), vlog.ClientProto(r.Proto), vlog.ClientURI(r.RequestURI), "Peers not found")
		return
	}

	peer, err := p.getNextPeer()
	if err != nil || peer == nil {
		http.Error(w, "Proxy error", http.StatusServiceUnavailable)
		p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(r.RemoteAddr), vlog.ClientHost(r.Host), vlog.ClientMethod(r.Method), vlog.ClientProto(r.Proto), vlog.ClientURI(r.RequestURI), err.Error())
		return
	}

	newProxyURI := fmt.Sprintf("%s://%s%s", peer.Proto, peer.Uri, r.RequestURI)

	newRequest, err := http.NewRequest(r.Method, newProxyURI, r.Body)
	if err != nil {
		http.Error(w, "Proxy error", http.StatusInternalServerError)
		p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(r.RemoteAddr), vlog.ClientHost(r.Host), vlog.ClientMethod(r.Method), vlog.ClientProto(r.Proto), vlog.ClientURI(r.RequestURI), err.Error())
		return
	}
	p.copyHeader(r.Header, &newRequest.Header)

	var transport http.Transport
	resp, err := transport.RoundTrip(newRequest)
	if err != nil {
		http.Error(w, "Proxy error", http.StatusInternalServerError)
		if resp != nil {
			p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(r.RemoteAddr), vlog.ClientHost(r.Host), vlog.ClientMethod(r.Method), vlog.ClientProto(r.Proto), vlog.ClientURI(r.RequestURI), vlog.ProxyHost(resp.Request.Host), vlog.ProxyMethod(resp.Request.Method), vlog.ProxyProto(resp.Request.Proto), vlog.ProxyURI(resp.Request.URL.Path), err.Error())
		} else {
			p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(r.RemoteAddr), vlog.ClientHost(r.Host), vlog.ClientMethod(r.Method),	vlog.ClientProto(r.Proto), vlog.ClientURI(r.RequestURI), err.Error())	
		}
		return
	}

	defer resp.Body.Close()
	resultBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Proxy error", http.StatusInternalServerError)
		p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(r.RemoteAddr), vlog.ClientHost(r.Host), vlog.ClientMethod(r.Method), vlog.ClientProto(r.Proto), vlog.ClientURI(r.RequestURI), vlog.ProxyHost(resp.Request.Host), vlog.ProxyMethod(resp.Request.Method), vlog.ProxyProto(resp.Request.Proto), vlog.ProxyURI(resp.Request.URL.Path), err.Error())
		return
	}

	_, err = w.Write(resultBody)
	if err != nil {
		p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(r.RemoteAddr), vlog.ClientHost(r.Host), vlog.ClientMethod(r.Method), vlog.ClientProto(r.Proto), vlog.ClientURI(r.RequestURI), vlog.ProxyHost(resp.Request.Host), vlog.ProxyMethod(resp.Request.Method), vlog.ProxyProto(resp.Request.Proto), vlog.ProxyURI(resp.Request.URL.Path), err.Error())
	}

	p.logger.Add(vlog.Debug, types.ResultOK, vlog.RemoteAddr(r.RemoteAddr), vlog.ClientHost(r.Host), vlog.ClientMethod(r.Method), vlog.ClientProto(r.Proto), vlog.ClientURI(r.RequestURI), vlog.ProxyHost(resp.Request.Host), vlog.ProxyMethod(resp.Request.Method), vlog.ProxyProto(resp.Request.Proto), vlog.ProxyURI(resp.Request.URL.Path))

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
