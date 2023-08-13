package proxy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"vbalancer/internal/proxy/peers"
	"vbalancer/internal/proxy/response"
	"vbalancer/internal/proxy/rules"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

const MaxCountCopyData = 2

var ErrCantGetProxyPort = errors.New("can't get proxy port")
var ErrMaxCountAttempts = errors.New("exceeded maximum number of attempts")
var ErrConfigPeersIsNil = errors.New("empty list peer in config file")

// Proxy defines the structure for the proxy server.
type Proxy struct {
	// 
	Logger vlog.ILog `yaml:"-" json:"-"`
	// Define the default port to listen on
	Port string `yaml:"port" json:"port"`
	// Define the client deadline time
	ClientDeadLineTime time.Duration `yaml:"clientDeadLineTime" json:"clientDeadLineTime"`
	// Define the peer host timeout
	PeerConnectionTimeout  time.Duration `yaml:"peerConnectionTimeout" json:"peerConnectionTimeout"`
	// Define the peer host deadline
	PeerHostDeadLine time.Duration `yaml:"peerHostDeadLine" json:"peerHostDeadLine"`
	// Define the max connection semaphore
	MaxCountConnection uint `yaml:"maxCountConnection" json:"maxCountConnection"`
	// Define the count dial attempts to peer
	CountMaxDialAttemptsToPeer uint `yaml:"countDialAttemptsToPeer" json:"countDialAttemptsToPeer"`
	// Peers is a list of peer configurations.
	Peers *peers.Peers `yaml:"peers" json:"peers"`
	// Defien allows configuration of blacklist rules to be passed to the proxy server
	Rules *rules.Rules `yaml:"rules" json:"rules"`
}

func New() *Proxy {
	return &Proxy{
		Logger:                     nil,
		Port:                       types.DefaultProxyPort,
		ClientDeadLineTime:         types.DeafultClientDeadLineTime,
		PeerConnectionTimeout:      types.DeafultPeerConnectionTimeout,
		PeerHostDeadLine:           types.DeafultPeerHostDeadLine,
		MaxCountConnection:         types.DeafultMaxCountConnection,
		CountMaxDialAttemptsToPeer: types.DeafultCountMaxDialAttemptsToPeer,
		//nolint:exhaustivestruct,exhaustruct
		Peers: &peers.Peers{},
		//nolint:exhaustivestruct,exhaustruct
		Rules: &rules.Rules{},
	}
}

// Init initializes the proxy server.
func (p *Proxy) Init(ctx context.Context, logger vlog.ILog) error {
	p.Logger = logger

	if p.Peers != nil && len(p.Peers.List) != 0 {
		err := p.Peers.Init(p.Peers.List)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	} else {
		return ErrConfigPeersIsNil
	}

	if p.Rules != nil {
		err := p.Rules.Init(ctx)
		if err != nil {
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

	defer func(proxySrv net.Listener) {
		err = proxySrv.Close()
		if err != nil {
			p.Logger.Add(vlog.Debug, types.ErrProxy, fmt.Errorf("proxy close failed: %w", err))
		}
	}(proxySrv)

	p.AcceptConnections(ctx, proxySrv)

	return nil
}

// AcceptConnections accepts connections from the proxy server.
func (p *Proxy) AcceptConnections(ctx context.Context, proxySrv net.Listener) {
	semaphore := make(chan struct{}, p.MaxCountConnection)

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

			go p.handleIncomingConnection(ctx, conn, semaphore)
		}
	}
}

// handleIncomingConnection. 
func (p *Proxy) handleIncomingConnection(ctx context.Context, conn net.Conn, semaphore chan struct{}) {
	defer func() {
		<-semaphore
	}()

	defer p.Logger.Add(vlog.Debug, types.ResultOK,
		fmt.Sprintf("accept connection %s, was close", conn.RemoteAddr().String()))

	defer conn.Close()

	p.Logger.Add(vlog.Debug, types.ResultOK, vlog.RemoteAddr(conn.RemoteAddr().String()), "starting connection")

	err := conn.SetDeadline(time.Now().Add(p.ClientDeadLineTime))
	if err != nil {
		p.Logger.Add(vlog.Debug, types.ErrProxy,
			vlog.RemoteAddr(conn.RemoteAddr().String()), fmt.Errorf("failed to set deadline: %w", err))

		return
	}

	clientAddr := conn.RemoteAddr().String()

	err = p.reverseData(ctx, conn, 0, p.CountMaxDialAttemptsToPeer)

	if err != nil {
		p.Logger.Add(vlog.Debug, types.ErrProxy, vlog.RemoteAddr(clientAddr), fmt.Errorf("failed in reverseData() %w", err))

		responseLogger := response.New()

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
func (p *Proxy) reverseData(ctx context.Context, client net.Conn,
	 curentDialAttemptsToPeer uint, maxDialAttemptsToPeer uint) error {
	if curentDialAttemptsToPeer >= maxDialAttemptsToPeer {
		return ErrMaxCountAttempts
	}

	pPeer, resultCode := p.Peers.GetNextPeer()
	if resultCode != types.ResultOK || pPeer == nil {
		//nolint:goerr113
		return fmt.Errorf("failed get next peer, result code: %s", resultCode.ToStr())
	}

	ctxConnectionTimeout, cancel := context.WithTimeout(ctx, p.PeerConnectionTimeout)
	defer cancel()	

	dst, err := pPeer.Dial(ctxConnectionTimeout, p.PeerConnectionTimeout, p.PeerHostDeadLine)
	if err != nil {
		curentDialAttemptsToPeer++

		return p.reverseData(ctx, client, curentDialAttemptsToPeer, maxDialAttemptsToPeer)
	}
	defer dst.Close()

	p.Logger.Add(
		vlog.Debug, 
		types.ResultOK,
		vlog.RemoteAddr(dst.RemoteAddr().String()),
		vlog.ProxyHost(client.LocalAddr().String()),
		fmt.Sprintf("try to copy data from remote: %s to peer: %s",
			client.RemoteAddr().String(), dst.RemoteAddr().String()))

	var waitGroup sync.WaitGroup

	waitGroup.Add(MaxCountCopyData)

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

// updatePort.
func (p *Proxy) updatePort() types.ResultCode {
	var proxyPort string

	if p.Port == "" || p.Port == ":" {
		proxyPort = os.Getenv("ProxyPort")
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
