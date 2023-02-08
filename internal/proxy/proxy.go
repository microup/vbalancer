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
			p.Logger.Add(types.Debug, types.ErrProxy, fmt.Sprintf("proxy close failed: %v\n", err))
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
					fmt.Sprintf("Accept failed, %v\n", err))
			} else {
				p.Logger.Add(types.Debug, types.ErrProxy, types.RemoteAddr("nil"),
					fmt.Sprintf("Accept failed, %v\n", err))
			}

			continue
		}

		semaphore <- struct{}{}

		err = conn.SetDeadline(time.Now().Add(time.Duration(p.Cfg.ClientDeadLineTimeSec) * time.Second))
		if err != nil {
			p.Logger.Add(types.Debug, types.ErrProxy, types.RemoteAddr(conn.RemoteAddr().String()),
				fmt.Sprintf("failed to set deadline: %v", err))
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
				p.Logger.Add(types.Debug, types.ResultOK, types.RemoteAddr(conn.RemoteAddr().String()), "starting connection")
				p.handleClientConnection(conn, 0, len(p.Peers.List))
				clientAddr := conn.RemoteAddr().String()
				err = conn.Close()

				if err != nil {
					p.Logger.Add(types.ErrCantFindActivePeers, types.ErrProxy, types.RemoteAddr(clientAddr),
						fmt.Errorf("failed client close %w", err))
				} else {
					p.Logger.Add(types.ErrCantFindActivePeers, types.ErrProxy, types.RemoteAddr(clientAddr),
						"the connection with the client was closed successfully")
				}
			}
		}(conn)
	}
}

func (p *Proxy) handleClientConnection(client net.Conn, numberOfAttempts int, maxNumberOfAttempts int) {
	if numberOfAttempts >= maxNumberOfAttempts {
		p.Logger.Add(
			types.ErrEmptyValue,
			types.ErrCantFindActivePeers,
			types.RemoteAddr(client.RemoteAddr().String()),
			types.ErrCantFindActivePeers)

		responseLogger := response.New(p.Logger)
		err := responseLogger.SentResponse(client, types.ErrProxy)

		if err != nil {
			p.Logger.Add(types.ErrCantFindActivePeers, types.ErrProxy, types.RemoteAddr(client.RemoteAddr().String()),
				fmt.Errorf("failed send response %w", err))
		}

		return
	}

	pPeer, resultCode := p.Peers.GetNextPeer()
	if resultCode != types.ResultOK || pPeer == nil {
		if pPeer != nil {
			p.Logger.Add(types.ErrCantFindActivePeers, resultCode, types.RemoteAddr(client.RemoteAddr().String()),
				types.ProxyHost(pPeer.GetURI()), resultCode.ToStr())
		} else {
			p.Logger.Add(types.ErrEmptyValue, resultCode, types.RemoteAddr(client.RemoteAddr().String()), resultCode.ToStr())
		}

		responseLogger := response.New(p.Logger)
		err := responseLogger.SentResponse(client, resultCode)

		if err != nil {
			p.Logger.Add(types.ErrCantFindActivePeers, types.ErrProxy, types.RemoteAddr(client.RemoteAddr().String()),
				fmt.Errorf("failed send response %w", err))
		}

		return
	}

	dst, err := net.DialTimeout("tcp", pPeer.GetURI(), time.Duration(p.Cfg.DestinationHostTimeoutMS)*time.Millisecond)
	if err != nil {
		numberOfAttempts++
		p.handleClientConnection(client, numberOfAttempts, maxNumberOfAttempts)

		return
	}

	err = dst.SetDeadline(time.Now().Add(time.Duration(p.Cfg.DestinationHostDeadLineSec) * time.Second))
	if err != nil {
		p.Logger.Add(types.Debug, types.ErrProxy, types.RemoteAddr(dst.LocalAddr().String()),
			fmt.Errorf("failed to set deadline: %w", err))

		return
	}

	p.ProxyDataCopy(client, dst)
}

//nolint:interfacer
func (p *Proxy) ProxyDataCopy(client net.Conn, dst net.Conn) {
	go func() {
		_, _ = bufio.NewReader(client).WriteTo(dst)
		dst.Close()
	}()

	_, _ = bufio.NewReader(dst).WriteTo(client)
	dst.Close()
}
