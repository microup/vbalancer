package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"vbalancer/internal/config"
	"vbalancer/internal/proxy"
	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

var ErrRecoveredPanic = errors.New("recovered from panic")

// Run this is the function of an application that starts a proxy server.
func Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cfg := config.New()

	err := cfg.Init()
	if err != nil {
		log.Panicf("failed to initialize configuration: %s", err.Error())
	}

	logger, err := vlog.New(cfg.Logger)
	if err != nil {
		log.Panicf("failed to create logger: %s", err)
	}

	defer func() {
		if err := recover(); err != nil {
			logger.Add(vlog.Fatal, types.ErrGotPanic, fmt.Errorf("%w: %v", ErrRecoveredPanic, err))
		}
	}()

	defer func(logger vlog.ILog) {
		err = logger.Close()
		if err != nil {
			log.Fatalf("failed close logger: %v", err)
		}
	}(logger)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Proxy.ShutdownTimeout)
	defer cancel()

	peerList := newPeerList(cfg)

	proxyBalancer := proxy.New(cfg.Proxy, cfg.Rules, peerList, logger)

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

	listenProxyChan := make(chan error)

	go func() {
		logger.Add(vlog.Info, types.ResultOK, fmt.Sprintf("start server addr on %s", cfg.ProxyPort))
		listenProxyChan <- proxyBalancer.ListenAndServe(ctx, cfg.ProxyPort)
	}()

	listenErr := <-listenProxyChan
	if listenErr != nil {
		logger.Add(vlog.Fatal, types.ErrProxy, fmt.Errorf("the proxy was return err: %w", err))
	}

	<-stopSignal
}

// newPeerList is the function that creates a list of peers for the balancer.
func newPeerList(cfg *config.Config) []peer.IPeer {
	listPeer := make([]peer.IPeer, len(cfg.Peers))

	for index, cfgPeer := range cfg.Peers {
		peerCopy := cfgPeer 
		listPeer[index] = &peerCopy
	}

	return listPeer
}
