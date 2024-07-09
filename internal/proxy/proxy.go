package proxy

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/proxy/peers"
	"vbalancer/internal/proxy/response"
	"vbalancer/internal/proxy/rules"

	"vbalancer/internal/types"
)

const MaxCountCopyData = 2
const BufferCopySize = 32 * 1024

var ErrCantGetProxyPort = errors.New("can't get proxy port")
var ErrMaxCountAttempts = errors.New("exceeded maximum number of attempts")
var ErrConfigPeersIsNil = errors.New("empty list peer in config file")

type Loger interface {
	Add(values ...interface{})
	Close() error
}

// Proxy defines the structure for the proxy server.
type Proxy struct {
	Logger Loger `json:"-" yaml:"-"`
	// Define the default port to listen on
	Port string `json:"port" yaml:"port"`
	// Define the client deadline time
	ClientDeadLineTime time.Duration `json:"clientDeadLineTime" yaml:"clientDeadLineTime"`
	// Define the peer host timeout
	PeerConnectionTimeout time.Duration `json:"peerConnectionTimeout" yaml:"peerConnectionTimeout"`
	// Define the max connection semaphore
	MaxCountConnection uint `json:"maxCountConnection" yaml:"maxCountConnection"`
	// Peers is a list of peer configurations.
	Peers *peers.Peers `json:"peers" yaml:"peers"`
	// Defien allows configuration of blacklist rules to be passed to the proxy server
	Rules *rules.Rules `json:"rules" yaml:"rules"`
}

func New() *Proxy {
	return &Proxy{
		Logger:                nil,
		Port:                  types.DefaultProxyPort,
		ClientDeadLineTime:    types.DeafultClientDeadLineTime,
		PeerConnectionTimeout: types.DeafultPeerConnectionTimeout,
		MaxCountConnection:    types.DeafultMaxCountConnection,
		Peers: &peers.Peers{
			TimeToEvictNotResponsePeers: 0,
			Peers:                       []peer.Peer{},
		},
		Rules: &rules.Rules{
			Blacklist: nil,
		},
	}
}

// Init initializes the proxy server.
func (p *Proxy) Init(ctx context.Context, logger Loger) error {
	p.Logger = logger

	if p.Peers != nil && len(p.Peers.Peers) != 0 {
		if err := p.Peers.Init(ctx, p.Peers.Peers); err != nil {
			return fmt.Errorf("%w", err)
		}
	} else {
		return ErrConfigPeersIsNil
	}

	if p.Rules != nil {
		if err := p.Rules.Init(ctx); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	if resultCode := p.updatePort(); resultCode != types.ResultOK {
		return fmt.Errorf("%w: %s", ErrCantGetProxyPort, resultCode.ToStr())
	}

	return nil
}

// ListenAndServe starts the proxy server.
func (p *Proxy) ListenAndServe(ctx context.Context) error {
	proxySrv, err := net.Listen("tcp", p.Port)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer proxySrv.Close()

	p.AcceptConnections(ctx, proxySrv)

	return nil
}

// AcceptConnections accepts connections from the proxy server.
func (p *Proxy) AcceptConnections(ctx context.Context, proxySrv net.Listener) {
	workers := make(chan struct{}, p.MaxCountConnection)
	defer close(workers)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := proxySrv.Accept()
			if err != nil {
				p.Logger.Add(types.Debug, types.ErrProxy, "remote accept connection failed", fmt.Errorf("%w", err))
				continue
			}

			if p.getCheckIsBlackListIP(conn.RemoteAddr().String()) {
				conn.Close()
				continue
			}

			if err = conn.SetDeadline(time.Now().Add(p.ClientDeadLineTime)); err != nil {
				p.Logger.Add(types.Debug, types.ErrProxy,
					types.RemoteAddr(conn.RemoteAddr().String()), "failed to set deadline", fmt.Errorf("%w", err))
				continue
			}

			workers <- struct{}{}

			go p.handleIncomingConnection(conn, workers)
		}
	}
}

// handleIncomingConnection.
func (p *Proxy) handleIncomingConnection(client net.Conn, worker chan struct{}) {
	defer func(logger Loger) {
		err := client.Close()
		if err != nil {
			logger.Add(types.Debug, types.ResultOK,
				types.RemoteAddr(client.LocalAddr().String()),
				types.PeerAddr(client.RemoteAddr().String()),
				"client closed with: ", err)
		} else {
			logger.Add(types.Error, types.ResultOK,
				types.RemoteAddr(client.LocalAddr().String()),
				types.PeerAddr(client.RemoteAddr().String()),
				"client closed OK")
		}

		<-worker
	}(p.Logger)

	if err := p.reverseData(client); err != nil {
		clientAddr := client.RemoteAddr().String()
		p.Logger.Add(types.Debug, types.ErrProxy, types.RemoteAddr(clientAddr),
			"failed in reverseData()", fmt.Errorf("%w", err))

		responseLogger := response.New()
		if err = responseLogger.SentResponseToClient(client, err); err != nil {
			p.Logger.Add(
				types.Debug,
				types.ErrSendResponseToClient,
				types.ErrProxy,
				types.RemoteAddr(clientAddr),
				"failed send response to client", fmt.Errorf("%w", err))
		}

		return
	}
}

// ReverseData reverses data from the client to the next available peer,
// it returns an error if the maximum number of attempts is reached or if it fails to get the next peer.
func (p *Proxy) reverseData(client net.Conn) error {
	pPeer, resultCode := p.Peers.GetNextPeer()
	if resultCode != types.ResultOK || pPeer == nil {
		return fmt.Errorf("failed get next peer, result code: %s", resultCode.ToStr())
	}

	dst, err := pPeer.Dial(p.PeerConnectionTimeout)
	if err != nil || dst == nil {
		p.Peers.AddToCacheBadPeer(pPeer.URI)

		return p.reverseData(client)
	}
	defer dst.Close()

	p.Logger.Add(
		types.Debug,
		types.ResultOK,
		types.RemoteAddr(client.RemoteAddr().String()),
		types.PeerAddr(dst.RemoteAddr().String()),
		"conneted to peer")

	p.proxyDataCopy(client, dst)

	p.Logger.Add(types.Debug, types.ResultOK,
		types.RemoteAddr(client.RemoteAddr().String()),
		types.PeerAddr(dst.RemoteAddr().String()),
		"copy data from peer to client was finish")

	return nil
}

// proxyDataCopy this is a function that copies data from the client to the peer
// and copies the response from the peer to the client.
func (p *Proxy) proxyDataCopy(client net.Conn, peer net.Conn) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(MaxCountCopyData)

	go func() {
		defer waitGroup.Done()
		p.copyData(client, peer)
	}()

	go func() {
		defer waitGroup.Done()
		p.copyData(peer, client)
	}()

	waitGroup.Wait()
}

// copyData copies data from src to dst until an error occurs or the deadline is reached.
func (p *Proxy) copyData(src net.Conn, dst net.Conn) {
	buffer := make([]byte, BufferCopySize)
	for {
		_ = src.SetReadDeadline(time.Now().Add(p.ClientDeadLineTime))
		n, err := src.Read(buffer)
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				return
			}
			return
		}
		_, err = dst.Write(buffer[:n])
		if err != nil {
			return
		}
	}
}

// updatePort.
func (p *Proxy) updatePort() types.ResultCode {
	var proxyPort string

	if p.Port == "" || p.Port == ":" {
		proxyPort = os.Getenv(types.ProxyPort)
		if proxyPort == ":" || proxyPort == "" {
			proxyPort = types.DefaultProxyPort
		}
	} else {
		proxyPort = p.Port
	}

	proxyPort = fmt.Sprintf(":%s", proxyPort)

	proxyPort = strings.Trim(proxyPort, " ")
	if proxyPort == strings.Trim(":", " ") {
		return types.ErrEmptyValue
	}

	p.Port = proxyPort

	return types.ResultOK
}

// getCheckIsBlackListIP.
func (p *Proxy) getCheckIsBlackListIP(remoteIP string) bool {
	if p.Rules != nil && p.Rules.Blacklist != nil {
		if p.Rules.Blacklist.IsBlacklistIP(remoteIP) {
			return true
		}
	}

	return false
}
