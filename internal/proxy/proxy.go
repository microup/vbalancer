package proxy

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"vbalancer/internal/config"
	"vbalancer/internal/proxy/peers"
	"vbalancer/internal/proxy/response"
	"vbalancer/internal/proxy/rules"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

const MaxProxyCountCopyData = 2

// Proxy defines the structure for the proxy server.
type Proxy struct {
	Logger vlog.ILog
	Peers  *peers.Peers
	Config *config.Proxy
	Rules  *rules.Rules
}

// New creates a new proxy server.
func New(cfg *config.Config, rules *rules.Rules, logger vlog.ILog) *Proxy {
	proxy := &Proxy{
		Config: cfg.Proxy,
		Logger: logger,
		Peers:  peers.New(cfg.Peers),
		Rules:  rules,
	}

	return proxy
}

// ListenAndServe starts the proxy server.
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

// AcceptConnections accepts connections from the proxy server.
func (p *Proxy) AcceptConnections(ctx context.Context, proxySrv net.Listener) {
	semaphore := make(chan struct{}, p.Config.MaxCountConnection)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := proxySrv.Accept()
			if err != nil {
				p.Logger.Add(vlog.Debug, types.ErrProxy, fmt.Errorf("accept connection failed: %w", err))

				continue
			}

			if p.getCheckIsBlackListIP(conn.RemoteAddr().String()) {
				conn.Close()

				continue
			}

			semaphore <- struct{}{}

			go p.handleIncomingConnection(conn, semaphore)
		}
	}
}

func (p *Proxy) getCheckIsBlackListIP(remoteIP string) bool {
	if p.Rules != nil && p.Rules.Blacklist != nil {
		if p.Rules.Blacklist.IsIPInBlacklist(remoteIP) {
			return true
		}
	}

	return false
}

func (p *Proxy) handleIncomingConnection(conn net.Conn, semaphore chan struct{}) {
	defer func() {
		<-semaphore
	}()

	defer p.Logger.Add(vlog.Debug, types.ResultOK,
		fmt.Sprintf("accept connection %s, was close", conn.RemoteAddr().String()))

	defer conn.Close()

	p.Logger.Add(vlog.Debug, types.ResultOK, vlog.RemoteAddr(conn.RemoteAddr().String()), "starting connection")

	err := conn.SetDeadline(time.Now().Add(p.Config.ClientDeadLineTime))
	if err != nil {
		p.Logger.Add(vlog.Debug, types.ErrProxy,
			vlog.RemoteAddr(conn.RemoteAddr().String()), fmt.Errorf("failed to set deadline: %w", err))

		return
	}

	clientAddr := conn.RemoteAddr().String()

	err = p.reverseData(conn, 0, p.Config.CountDialAttemptsToPeer)

	if err != nil {
		p.Logger.Add(vlog.Debug, types.ErrProxy, vlog.RemoteAddr(clientAddr), fmt.Errorf("failed in reverseData() %w", err))

		responseLogger := response.New(p.Logger)

		err = responseLogger.SentResponseToClient(conn, err)
		if err != nil {
			p.Logger.Add(
				vlog.Debug,
				types.ErrSendResponseToClient,
				types.ErrProxy,
				vlog.RemoteAddr(clientAddr),
				fmt.Errorf("failed send response to client %w", err))
		}
	}
}

// ReverseData reverses data from the client to the next available peer,
// it returns an error if the maximum number of attempts is reached or if it fails to get the next peer.
func (p *Proxy) reverseData(client net.Conn, numberOfAttempts uint, maxNumberOfAttempts uint) error {
	if numberOfAttempts >= maxNumberOfAttempts {
		return types.ErrMaxCountAttempts
	}

	pPeer, resultCode := p.Peers.GetNextPeer()
	if resultCode != types.ResultOK || pPeer == nil {
		//nolint:goerr113
		return fmt.Errorf("failed get next peer, result code: %s", resultCode.ToStr())
	}

	dst, err := pPeer.Dial(p.Config.PeerHostTimeOut, p.Config.PeerHostDeadLine)
	if err != nil {
		numberOfAttempts++

		return p.reverseData(client, numberOfAttempts, maxNumberOfAttempts)
	}
	defer dst.Close()

	p.Logger.Add(vlog.Debug, types.ResultOK,
		vlog.RemoteAddr(dst.RemoteAddr().String()),
		vlog.ProxyHost(client.LocalAddr().String()),
		fmt.Sprintf("try to copy data from remote: %s to peer: %s",
			client.RemoteAddr().String(), dst.RemoteAddr().String()))

	var waitGroup sync.WaitGroup

	waitGroup.Add(MaxProxyCountCopyData)

	p.proxyDataCopy(&waitGroup, client, dst)

	waitGroup.Wait()

	p.Logger.Add(vlog.Debug, types.ResultOK,
		vlog.RemoteAddr(dst.RemoteAddr().String()),
		vlog.ProxyHost(client.LocalAddr().String()),
		fmt.Sprintf("copy data from remote: %s to peer: %s, it was finish",
			client.RemoteAddr().String(), dst.RemoteAddr().String()))

	return nil
}

// proxyDataCopy this is a function that copies data from the client to the peer
// and copies the response from the peer to the client.
func (p *Proxy) proxyDataCopy(waitGroup *sync.WaitGroup, client io.ReadWriter, dst io.ReadWriter) {
	go func() {
		defer waitGroup.Done()

		_, _ = io.Copy(dst, client)
	}()

	go func() {
		defer waitGroup.Done()

		_, _ = io.Copy(client, dst)
	}()
}
