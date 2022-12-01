package proxy

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"
	"vbalancer/internal/peer"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

const (
	maxCopyChan = 2
)

type Proxy struct {
	srv              net.Listener
	logger           *vlog.VLog
	Peers            []*peer.Peer
	cfg              *Config
	CurrentPeerIndex *uint64
}

func New(proxyPort string, cfg *Config, peers []*peer.Peer, logger *vlog.VLog) *Proxy {
	proxySrv, err := net.Listen("tcp", proxyPort)
	if err != nil {
		panic("connection error:" + err.Error())
	}

	var startPeerIndex uint64

	proxy := &Proxy{
		srv:              proxySrv,
		logger:           logger,
		Peers:            peers,
		cfg:              cfg,
		CurrentPeerIndex: &startPeerIndex,
	}

	return proxy
}

func (p *Proxy) Start(ctx context.Context, checkTimeAlive *peer.CheckTimeAlive) error {
	for _, peer := range p.Peers {
		peer.CheckTimeAlive = checkTimeAlive
		peer.Logger = p.logger

		go peer.CheckIsAlive(ctx)
	}

	for {
		if conn, err := p.srv.Accept(); err == nil {
			err = conn.SetDeadline(time.Now().Add(time.Duration(p.cfg.TimeDeadLineMS) * time.Millisecond))
			if err != nil {
				p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(conn.RemoteAddr().String()),
					fmt.Sprintf("Accept failed, %v\n", err))
			}

			go p.copyConn(conn)
		} else {
			p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(conn.RemoteAddr().String()),
				fmt.Sprintf("Accept failed, %v\n", err))
		}
	}
}

func (p *Proxy) Shutdown(ctx context.Context) error {
	return nil
}

func (p *Proxy) copyConn(client net.Conn) {
	defer client.Close()

	peer, resultCode := p.GetNextPeer()

	if resultCode != types.ResultOK || peer == nil {
		if peer != nil {
			p.logger.Add(vlog.Debug, resultCode, vlog.RemoteAddr(client.RemoteAddr().String()),
				vlog.ProxyHost(peer.URI))
		} else {
			p.logger.Add(vlog.Debug, resultCode, vlog.RemoteAddr(client.RemoteAddr().String()))
		}

		return
	}

	dst, err := net.Dial("tcp", peer.URI)
	if err != nil {
		p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(client.RemoteAddr().String()),
			fmt.Sprintf("Accept failed, %v\n", err))

		return
	}

	p.logger.Add(vlog.Debug,
		types.ResultOK,
		vlog.RemoteAddr(client.RemoteAddr().String()),
		vlog.ProxyHost(peer.URI),
	)

	done := make(chan bool, maxCopyChan)

	go func() {
		defer client.Close()
		defer dst.Close()
		//nolint:errcheck
		io.Copy(dst, client)
		done <- true
	}()

	go func() {
		defer client.Close()
		defer dst.Close()
		//nolint:errcheck
		io.Copy(client, dst)
		done <- true
	}()

	<-done
	<-done
}

func (p *Proxy) GetNextPeer() (*peer.Peer, types.ResultCode) {
	var next int

	if *p.CurrentPeerIndex >= uint64(len(p.Peers)) {
		atomic.StoreUint64(p.CurrentPeerIndex, uint64(0))
	} else {
		next = p.nextIndex()
	}

	l := len(p.Peers) + next
	for i := next; i < l; i++ {
		idx := i % len(p.Peers)
		isAlive := p.Peers[idx].IsAlive()

		if isAlive {
			atomic.StoreUint64(p.CurrentPeerIndex, uint64(idx))

			return p.Peers[idx], types.ResultOK
		}
	}

	return nil, types.ErrCantFinePeers
}

func (p *Proxy) nextIndex() int {
	return int(atomic.AddUint64(p.CurrentPeerIndex, uint64(1)) % uint64(len(p.Peers)))
}
