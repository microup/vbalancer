package proxy

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
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
	Logger vlog.ILog
	Peers  *peers.Peers
	Cfg    *Config
}

func New(cfg *Config, listPeer []peer.IPeer, logger vlog.ILog) *Proxy {
	proxy := &Proxy{
		Logger: logger,
		Peers:  peers.New(listPeer),
		Cfg:    cfg,
	}

	return proxy
}

func (p *Proxy) Start(ctx context.Context, proxyPort string, checkTimeAlive *peer.CheckTimeAlive) error {
	proxySrv, err := net.Listen("tcp", proxyPort)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer func(proxySrv net.Listener) {
		err = proxySrv.Close()
		if err != nil {
			p.Logger.Add(vlog.Debug, types.ErrProxy, fmt.Sprintf("proxy close failed: %v\n", err))
		}
	}(proxySrv)

	for _, pPeer := range p.Peers.List {
		pPeer.SetCheckTimeAlive(checkTimeAlive)
		pPeer.SetLogger(p.Logger)

		go pPeer.CheckIsAlive(ctx)
	}

	p.checkNewConnection(proxySrv)

	return nil
}

func (p *Proxy) checkNewConnection(proxySrv net.Listener) {
	for {
		conn, err := proxySrv.Accept()
		if err != nil {
			p.Logger.Add(vlog.Debug, types.ErrProxy, vlog.RemoteAddr(conn.RemoteAddr().String()),
				fmt.Sprintf("Accept failed, %v\n", err))

			continue
		}

		err = conn.SetDeadline(time.Now().Add(time.Duration(p.Cfg.DeadLineTimeMS) * time.Millisecond))
		if err != nil {
			p.Logger.Add(vlog.Debug, types.ErrProxy, vlog.RemoteAddr(conn.RemoteAddr().String()),
				fmt.Sprintf("failed to set deadline: %v", err))

			continue
		}

		go p.handleConnection(conn)
	}
}

func (p *Proxy) Shutdown() error {
	return nil
}

func (p *Proxy) handleConnection(client net.Conn) {
	defer func(client net.Conn) {
		_ = client.Close()
	}(client)

	p.Logger.Add(vlog.Debug, types.ResultOK, vlog.RemoteAddr(client.RemoteAddr().String()))

	pPeer, resultCode := p.Peers.GetNextPeer()

	if resultCode != types.ResultOK || pPeer == nil {
		if pPeer != nil {
			p.Logger.Add(vlog.Debug, resultCode, vlog.RemoteAddr(client.RemoteAddr().String()),
				vlog.ProxyHost(pPeer.GetURI()), resultCode.ToStr())
		} else {
			p.Logger.Add(vlog.Debug, resultCode, vlog.RemoteAddr(client.RemoteAddr().String()), resultCode.ToStr())
		}

		responseLogger := response.New(p.Logger)
		_ = responseLogger.SentResponse(client, resultCode)

		return
	}

	dst, err := net.DialTimeout("tcp", pPeer.GetURI(), time.Duration(p.Cfg.DeadLineTimeMS)*time.Millisecond)
	if err != nil {
		p.Logger.Add(vlog.Debug, types.ErrProxy, vlog.RemoteAddr(client.RemoteAddr().String()),
			fmt.Sprintf("failed connecting to target:, %v\n", err))

		responseLogger := response.New(p.Logger)
		_ = responseLogger.SentResponse(client, types.ErrProxy)

		return
	}

	p.Logger.Add(vlog.Debug, types.ResultOK, vlog.RemoteAddr(client.RemoteAddr().String()),
		vlog.ProxyHost(pPeer.GetURI()))

	p.ProxyDataCopy(client, dst)
}

func (p *Proxy) ProxyDataCopy(client net.Conn, dst io.ReadWriteCloser) {
	done := make(chan struct{}, maxCopyChan)
	defer close(done)

	p.CopyDataPeerToClient(dst, client, done)
	p.copyDataClientToPeer(client, dst, done)

	<-done
	<-done
}

func (p *Proxy) copyDataClientToPeer(client net.Conn, dst io.ReadCloser, done chan struct{}) {
	go func() {
		writeBuffer := make([]byte, p.Cfg.SizeCopyBufferIO)
		_, _ = io.CopyBuffer(client, dst, writeBuffer)

		_ = dst.Close()
		_ = client.Close()
		done <- struct{}{}
	}()
}

func (p *Proxy) CopyDataPeerToClient(dst io.WriteCloser, client net.Conn, done chan struct{}) {
	go func() {
		readBuffer := make([]byte, p.Cfg.SizeCopyBufferIO)
		_, _ = io.CopyBuffer(dst, client, readBuffer)

		_ = dst.Close()
		_ = client.Close()
		done <- struct{}{}
	}()
}
