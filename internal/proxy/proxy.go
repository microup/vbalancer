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
			p.Logger.Add(vlog.Debug, types.ErrProxy, fmt.Errorf("proxy close failed: %w", err))
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
				p.Logger.Add(vlog.Debug, types.ErrProxy, vlog.RemoteAddr(conn.RemoteAddr().String()),
					fmt.Errorf("accept failed, %w", err))
			} else {
				p.Logger.Add(vlog.Debug, types.ErrProxy, vlog.RemoteAddr("nil"),
					fmt.Errorf("accept failed, %w", err))
			}

			continue
		}

		semaphore <- struct{}{}

		err = conn.SetDeadline(time.Now().Add(p.Cfg.ClientDeadLineTime))
		if err != nil {
			p.Logger.Add(vlog.Debug, types.ErrProxy, vlog.RemoteAddr(conn.RemoteAddr().String()),
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

	p.Logger.Add(vlog.Debug, types.ResultOK,
		vlog.RemoteAddr(conn.RemoteAddr().String()),
		"starting connection")

	err := p.reverseData(conn, 0, p.Cfg.CountDialAttemptsToPeer)
	if err != nil {
		p.Logger.Add(vlog.Debug, types.ErrProxy, vlog.RemoteAddr(clientAddr),
			fmt.Errorf("failed in reverseData() %w", err))

		responseLogger := response.New(p.Logger)
		
		err = responseLogger.SentResponseToClient(conn, err)

		if err != nil {
			p.Logger.Add(vlog.Debug, types.ErrSendResponseToClient, types.ErrProxy,
				vlog.RemoteAddr(clientAddr),
				fmt.Errorf("failed send response to client %w", err))
		}
	}

	err = conn.Close()

	if err != nil {
		p.Logger.Add(vlog.Debug, types.ErrProxy, types.ErrProxy,
			vlog.RemoteAddr(clientAddr),
			fmt.Errorf("failed client close %w", err))
	} else {
		p.Logger.Add(vlog.Debug, types.ResultOK, vlog.RemoteAddr(clientAddr),
			"the connection with the client was closed successfully")
	}
}

// ReverseData - reverses data from the client to the next available peer.
// It returns an error if the maximum number of attempts is reached or if it fails to get the next peer.
func (p *Proxy) reverseData(client net.Conn, numberOfAttempts uint, maxNumberOfAttempts uint) error {
	if numberOfAttempts >= maxNumberOfAttempts {
		return types.ErrMaxCountAttempts
	}

	pPeer, resultCode := p.Peers.GetNextPeer()
	if resultCode != types.ResultOK || pPeer == nil {
		//nolint:goerr113
		return fmt.Errorf("failed get next peer, result code: %s", resultCode.ToStr())
	}

	dst, err := pPeer.Dial(p.Cfg.DestinationHostTimeOut, p.Cfg.DestinationHostDeadLine)
	if err != nil {
		numberOfAttempts++

		return p.reverseData(client, numberOfAttempts, maxNumberOfAttempts)
	}
	defer dst.Close()

	p.proxyDataCopy(client, dst)

	return nil
}

// proxyDataCopy - this is a function that copies data from the client to the peer
// and copies the response from the peer to the client.
func (p *Proxy) proxyDataCopy(client net.Conn, dst net.Conn) {
	p.Logger.Add(vlog.Debug, types.ResultOK,
		vlog.RemoteAddr(dst.RemoteAddr().String()),
		vlog.ProxyHost(client.LocalAddr().String()), "try to send data")

		go func() {
		_, _ = bufio.NewReader(client).WriteTo(dst)
	}()

	_, _ = bufio.NewReader(dst).WriteTo(client)
}
