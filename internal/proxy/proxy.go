package proxy

import (
	"bufio"
	"context"
	"fmt"
	"net"

	"time"
	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/proxy/peers"
	"vbalancer/internal/proxy/response"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
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

func (p *Proxy) ListenAndServe(ctx context.Context, proxyPort string) error {
	proxySrv, err := net.Listen("tcp", proxyPort)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer func(proxySrv net.Listener) {
		err = proxySrv.Close()
		if err != nil {
			p.Logger.Add(types.Debug, types.ErrProxy, fmt.Errorf("proxy close failed: %w", err))
		}
	}(proxySrv)

	for _, pPeer := range p.Peers.List {
		pPeer.SetLogger(p.Logger)
	}

	p.AcceptConnections(ctx, proxySrv)

	return nil
}

func (p *Proxy) AcceptConnections(ctx context.Context, proxySrv net.Listener) {
	semaphore := make(chan struct{}, p.Cfg.ConnectionSemaphore)

	for {
		conn, err := proxySrv.Accept()
		if err != nil {
			if conn != nil {
				p.Logger.Add(types.Debug, types.ErrProxy, types.RemoteAddr(conn.RemoteAddr().String()),
					fmt.Errorf("accept failed, %w", err))
			} else {
				p.Logger.Add(types.Debug, types.ErrProxy, types.RemoteAddr("nil"),
					fmt.Errorf("accept failed, %w", err))
			}

			continue
		}

		semaphore <- struct{}{}

		err = conn.SetDeadline(time.Now().Add(time.Duration(p.Cfg.ClientDeadLineTimeSec) * time.Second))
		if err != nil {
			p.Logger.Add(types.Debug, types.ErrProxy, types.RemoteAddr(conn.RemoteAddr().String()),
				fmt.Errorf("failed to set deadline: %w", err))
			<-semaphore

			continue
		}

		go func(conn net.Conn) {
			defer func() {
				<-semaphore
			}()

			select {
			case <-ctx.Done():
				return
			default:
				p.executeConnection(conn)
			}
		}(conn)
	}
}

func (p *Proxy) executeConnection(conn net.Conn) {
	clientAddr := conn.RemoteAddr().String()

	p.Logger.Add(types.Debug, types.ResultOK,
		types.RemoteAddr(conn.RemoteAddr().String()),
		"starting connection")

	err := p.reverseData(conn, 0, p.Cfg.CountDialAttemptsToPeer)
	if err != nil {
		p.Logger.Add(types.Debug, types.ErrProxy, types.RemoteAddr(clientAddr),
			"failed in reverseData() %w", err)

		responseLogger := response.New(p.Logger)
		
		err = responseLogger.SentResponseToClient(conn, err)

		if err != nil {
			p.Logger.Add(types.Debug, types.ErrSendResponseToClient, types.ErrProxy,
				types.RemoteAddr(clientAddr),
				fmt.Errorf("failed send response to client %w", err))
		}
	}

	err = conn.Close()

	if err != nil {
		p.Logger.Add(types.Debug, types.ErrProxy, types.ErrProxy,
			types.RemoteAddr(clientAddr),
			fmt.Errorf("failed client close %w", err))
	} else {
		p.Logger.Add(types.Debug, types.ErrProxy, types.RemoteAddr(clientAddr),
			"the connection with the client was closed successfully")
	}
}

func (p *Proxy) reverseData(client net.Conn, numberOfAttempts uint, maxNumberOfAttempts uint) error {
	if numberOfAttempts >= maxNumberOfAttempts {
		return types.ErrMaxCountAttempts
	}

	pPeer, resultCode := p.Peers.GetNextPeer()
	if resultCode != types.ResultOK || pPeer == nil {
		//nolint:goerr113
		return fmt.Errorf("failed get next peer, result code: %s", resultCode.ToStr())
	}

	dst, err := pPeer.Dial(
		time.Duration(p.Cfg.DestinationHostTimeOutMs)*time.Millisecond,
		time.Duration(p.Cfg.DestinationHostDeadLineSec)*time.Second)
	if err != nil {
		numberOfAttempts++

		return p.reverseData(client, numberOfAttempts, maxNumberOfAttempts)
	}
	defer dst.Close()

	p.proxyDataCopy(client, dst)

	return nil
}

func (p *Proxy) proxyDataCopy(client net.Conn, dst net.Conn) {
	p.Logger.Add(types.Debug, types.ResultOK,
		types.RemoteAddr(dst.RemoteAddr().String()),
		types.ProxyHost(client.LocalAddr().String()), "try to send data")

	go func() {
		_, _ = bufio.NewReader(client).WriteTo(dst)
	}()

	_, _ = bufio.NewReader(dst).WriteTo(client)
}
