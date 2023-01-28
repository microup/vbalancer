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
	chunkSize   = 32 * 1024
)

type Proxy struct {
	logger vlog.ILog
	Peers  *peers.Peers
	cfg    *Config
}

func New(cfg *Config, listPeer []peer.IPeer, logger vlog.ILog) *Proxy {
	proxy := &Proxy{
		logger: logger,
		Peers:  peers.New(listPeer),
		cfg:    cfg,
	}

	return proxy
}

func (p *Proxy) Start(ctx context.Context, proxyPort string, checkTimeAlive *peer.CheckTimeAlive) error {
	proxySrv, err := net.Listen("tcp", proxyPort)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for _, pPeer := range p.Peers.List {
		pPeer.SetCheckTimeAlive(checkTimeAlive)
		pPeer.SetLogger(p.logger)

		go pPeer.CheckIsAlive(ctx)
	}

	return p.checkNewConnection(proxySrv)
}

func (p *Proxy) checkNewConnection(proxySrv net.Listener) error {
	for {
		conn, err := proxySrv.Accept()
		if err != nil {
			p.logger.Add(vlog.Debug, types.ErrProxy, vlog.RemoteAddr(conn.RemoteAddr().String()),
				fmt.Sprintf("Accept failed, %v\n", err))

			continue
		}

		err = conn.SetDeadline(time.Now().Add(time.Duration(p.cfg.TimeDeadLineMS) * time.Millisecond))
		if err != nil {
			p.logger.Add(vlog.Debug, types.ErrProxy, vlog.RemoteAddr(conn.RemoteAddr().String()),
				fmt.Sprintf("Accept failed: %v\n", err))

			return fmt.Errorf("failed to set deadline: %w", err)
		}

		go p.copyConn(conn)
	}
}

func (p *Proxy) Shutdown() error {
	return nil
}

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

	p.proxyDataCopy(client, dst)

}

func (p *Proxy) proxyDataCopy(client net.Conn, dst io.ReadWriteCloser) {
	done := make(chan struct{}, maxCopyChan)
	defer close(done)

	p.copyDataClientToPeer(client, dst, done)
	p.copyDataPeerToClient(dst, client, done)

	<-done
	<-done

	_ = dst.Close()
}

func (p *Proxy) copyDataClientToPeer(client net.Conn, dst io.ReadWriteCloser, done chan struct{}) {
	writeBuffer := make([]byte, chunkSize)

	go func() {
		var errCopy error
		_, errCopy = io.CopyBuffer(client, dst, writeBuffer)

		if errCopy != nil {
			p.logger.Add(vlog.Debug, types.ErrCopyDataClientToPeer, vlog.RemoteAddr(client.RemoteAddr().String()),
				fmt.Sprintf(types.ErrorCopyDataClientToPeer, errCopy))
		}

		_ = client.Close()
		done <- struct{}{}
	}()
}

func (p *Proxy) copyDataPeerToClient(dst io.ReadWriteCloser, client net.Conn, done chan struct{}) {
	readBuffer := make([]byte, chunkSize)

	go func() {
		var errCopy error
		_, errCopy = io.CopyBuffer(dst, client, readBuffer)

		if errCopy != nil {
			p.logger.Add(vlog.Debug, types.ErrCopyDataPeerToClient, vlog.RemoteAddr(client.RemoteAddr().String()),
				fmt.Sprintf(types.ErrorCopyDataPeerToClient, errCopy))
		}

		_ = client.Close()
		done <- struct{}{}
	}()
}
