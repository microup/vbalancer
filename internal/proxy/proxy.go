package proxy

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"

	//"sync"
	"time"
	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/proxy/peers"
	"vbalancer/internal/proxy/response"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

const (
	maxCopyChannel = 2
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

func (p *Proxy) ListenAndServe(ctx context.Context, proxyPort string, checkTimeAlive *peer.CheckTimeAlive) error {
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
		pPeer.SetAvailabilityCheckInterval(checkTimeAlive)
		pPeer.SetLogger(p.Logger)

		go pPeer.CheckAvailability(ctx)
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

		err = conn.SetDeadline(time.Now().Add(time.Duration(p.Cfg.DeadLineTimeSeconds) * time.Second))
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
				p.handleClientConnection(conn)
			}
		}(conn)
	}
}

func (p *Proxy) handleClientConnection(client net.Conn) {

	defer func(client net.Conn) {
		err := client.Close()
		p.Logger.Add(types.ErrCantFindActivePeers, types.ErrProxy, types.RemoteAddr(client.RemoteAddr().String()),
			fmt.Sprintf("failed client close %v\n", err))
	}(client)

	p.Logger.Add(types.Debug, types.ResultOK, types.RemoteAddr(client.RemoteAddr().String()))

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
		p.Logger.Add(types.ErrCantFindActivePeers, types.ErrProxy, types.RemoteAddr(client.RemoteAddr().String()),
			fmt.Sprintf("failed send response %v\n", err))

		return
	}

	dst, err := net.DialTimeout("tcp", pPeer.GetURI(), time.Duration(p.Cfg.DeadLineTimeSeconds)*time.Millisecond)
	if err != nil {
		p.Logger.Add(types.Debug, types.ErrProxy, types.RemoteAddr(client.RemoteAddr().String()),
			fmt.Sprintf("failed connecting to target:, %v\n", err))

		responseLogger := response.New(p.Logger)
		err = responseLogger.SentResponse(client, types.ErrProxy)
		p.Logger.Add(types.ErrCantFindActivePeers, types.ErrProxy, types.RemoteAddr(client.RemoteAddr().String()),
			fmt.Sprintf("failed send response %v\n", err))

		return
	}

	p.Logger.Add(types.Debug, types.ResultOK, types.RemoteAddr(client.RemoteAddr().String()),
		types.ProxyHost(pPeer.GetURI()))

	p.ProxyDataCopy(client, dst)
}

func (p *Proxy) ProxyDataCopy(client net.Conn, dst io.ReadWriteCloser) {
	defer client.Close()

	//waitG := &sync.WaitGroup{}
	//waitG.Add(maxCopyChannel)

	//go p.copyClientToPeer(client, dst, waitG)
	//go p.copyPeerToClient(dst, client, waitG)

	go func() {
		_, err := bufio.NewReader(client).WriteTo(dst)
		if err != nil {
			fmt.Println("Error forwarding data from client to remote host:", err)
		}
		dst.Close()
	}()

	_, err := bufio.NewReader(dst).WriteTo(client)
	if err != nil {
		fmt.Println("Error forwarding data from remote host to client:", err)
	}

	dst.Close()
	//waitG.Wait()
}

// func (p *Proxy) copyClientToPeer(client net.Conn, dst io.ReadCloser, waitG *sync.WaitGroup) {
// 	defer func() {
// 		dst.Close()
// 		client.Close()
// 		waitG.Done()
// 	}()

// 	writeBuffer := make([]byte, p.Cfg.SizeCopyBufferIO)
// 	_, _ = io.CopyBuffer(client, dst, writeBuffer)
// }

// func (p *Proxy) copyPeerToClient(dst io.WriteCloser, client net.Conn, waitG *sync.WaitGroup) {
// 	defer func() {
// 		dst.Close()
// 		client.Close()
// 		waitG.Done()
// 	}()

// 	readBuffer := make([]byte, p.Cfg.SizeCopyBufferIO)
// 	_, _ = io.CopyBuffer(dst, client, readBuffer)
// }
