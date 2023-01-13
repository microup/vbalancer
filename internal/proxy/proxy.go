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

func New(proxyPort string, cfg *Config, listPeer []peer.IPeer, logger *vlog.VLog) *Proxy {
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
	for _, pPeer := range p.Peers.List {
		pPeer.SetCheckTimeAlive(checkTimeAlive)
		pPeer.SetLogger(p.logger)

		go pPeer.CheckIsAlive(ctx)
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

func (p *Proxy) Shutdown() error {
	return nil
}

//nolint:funlen
func (p *Proxy) copyConn(client net.Conn) {
	defer func(client net.Conn) {
		_ = client.Close()
	}(client)

	p.logger.Add(vlog.Debug, types.ResultOK, vlog.RemoteAddr(client.RemoteAddr().String()))

	pPeer, resultCode := p.Peers.GetNextPeer()

	if resultCode != types.ResultOK || pPeer == nil {
		if pPeer != nil {
			p.logger.Add(vlog.Debug, resultCode, vlog.RemoteAddr(client.RemoteAddr().String()),
				vlog.ProxyHost(pPeer.GetURI()), resultCode.ToStr())
		} else {
			p.logger.Add(vlog.Debug, resultCode, vlog.RemoteAddr(client.RemoteAddr().String()), resultCode.ToStr())
		}

		responseLogger := response.New(p.logger)
		_ = responseLogger.SentResponse(client, resultCode)

		return
	}

	dst, err := net.Dial("tcp", pPeer.GetURI())
	if err != nil {
		p.logger.Add(vlog.Debug, types.ErrProxy, vlog.RemoteAddr(client.RemoteAddr().String()),
			fmt.Sprintf("Accept failed, %v\n", err))

		responseLogger := response.New(p.logger)
		_ = responseLogger.SentResponse(client, types.ErrProxy)

		return
	}

	p.logger.Add(vlog.Debug, types.ResultOK, vlog.RemoteAddr(client.RemoteAddr().String()),
		vlog.ProxyHost(pPeer.GetURI()))

	done := make(chan bool, maxCopyChan)
	defer close(done)

	go func() {
		defer func(client net.Conn) {
			_ = client.Close()
		}(client)
		defer func(dst net.Conn) {
			_ = dst.Close()
		}(dst)

		_, _ = io.Copy(dst, client)
		done <- true
	}()

	go func() {
		defer func(client net.Conn) {
			_ = client.Close()
		}(client)
		defer func(dst net.Conn) {
			_ = dst.Close()
		}(dst)

		_, _ = io.Copy(client, dst)
		done <- true
	}()

	<-done
	<-done
}
