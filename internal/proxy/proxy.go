package proxy

import (
	"context"
	"time"

	"fmt"
	"io"
	"net"
	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/proxy/peers"
	"vbalancer/internal/proxy/response"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

const (
	maxCopyChan = 2
)

type Proxy struct {
	srv    net.Listener
	logger *vlog.VLog
	Peers  *peers.Peers
	cfg    *Config
}

func New(proxyPort string, cfg *Config, listPeer []*peer.Peer, logger *vlog.VLog) *Proxy {
	proxySrv, err := net.Listen("tcp", proxyPort)
	if err != nil {
		panic("connection error:" + err.Error())
	}

	proxy := &Proxy{
		srv:    proxySrv,
		logger: logger,
		Peers:  peers.New(listPeer),
		cfg:    cfg,
	}

	return proxy
}

func (p *Proxy) Start(ctx context.Context, checkTimeAlive *peer.CheckTimeAlive) error {
	for _, peer := range p.Peers.List {
		peer.CheckTimeAlive = checkTimeAlive
		peer.Logger = p.logger

		go peer.CheckIsAlive(ctx)
	}

	for {
		if conn, err := p.srv.Accept(); err == nil {
			err = conn.SetDeadline(time.Now().Add(time.Duration(p.cfg.TimeDeadLineMS) * time.Millisecond))

			if err != nil {
				p.logger.Add(vlog.Debug, types.ErrProxy, vlog.RemoteAddr(conn.RemoteAddr().String()),
					fmt.Sprintf("Accept failed, %v\n", err))
			}

			go p.copyConn(conn)
		} else {
			p.logger.Add(vlog.Debug, types.ErrProxy, vlog.RemoteAddr(conn.RemoteAddr().String()),
				fmt.Sprintf("Accept failed, %v\n", err))
		}
	}
}

func (p *Proxy) Shutdown(ctx context.Context) error {
	return nil
}

func (p *Proxy) copyConn(client net.Conn) {
	defer client.Close()

	p.logger.Add(vlog.Debug, types.ResultOK, vlog.RemoteAddr(client.RemoteAddr().String()))

	peer, resultCode := p.Peers.GetNextPeer()

	if resultCode != types.ResultOK || peer == nil {
		if peer != nil {
			p.logger.Add(vlog.Debug, resultCode, vlog.RemoteAddr(client.RemoteAddr().String()),
				vlog.ProxyHost(peer.URI), resultCode.ToStr())
		} else {
			p.logger.Add(vlog.Debug, resultCode, vlog.RemoteAddr(client.RemoteAddr().String()), resultCode.ToStr())
		}

		response := response.New(p.logger)
		response.SentResponse(client, resultCode)

		return
	}

	dst, err := net.Dial("tcp", peer.URI)
	if err != nil {
		p.logger.Add(vlog.Debug, types.ErrProxy, vlog.RemoteAddr(client.RemoteAddr().String()),
			fmt.Sprintf("Accept failed, %v\n", err))

		response := response.New(p.logger)
		response.SentResponse(client, types.ErrProxy)

		return
	}

	p.logger.Add(vlog.Debug,
		types.ResultOK,
		vlog.RemoteAddr(client.RemoteAddr().String()),
		vlog.ProxyHost(peer.URI))

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
