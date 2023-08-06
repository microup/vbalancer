package app

import (
	"context"
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

// Run this is the function of an application that starts a proxy server.
func Run() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("catch err: %v", err) //nolint:forbidigo
		}
	}()

	configuration := initializeConfig()

	logger, err := vlog.New(configuration.Logger)
	if err != nil {
		log.Panicf("failed created logger: %v", err)
	}

	defer func(logger vlog.ILog) {
		err = logger.Close()
		if err != nil {
			log.Fatalf("failed closed logger: %v", err)
		}
	}(logger)

	ctx, cancel := context.WithTimeout(context.Background(), configuration.Proxy.ShutdownTimeout)
	defer cancel()
	
	peerList := createPeerListForBalancer(configuration)

	proxyBalancer := proxy.New(configuration.Proxy, configuration.Rules, peerList, logger)

	logger.Add(vlog.Info, types.ResultOK, fmt.Sprintf("start server addr on %s", configuration.ProxyPort))

	go func() {
		if err = proxyBalancer.ListenAndServe(ctx, configuration.ProxyPort); err != nil {
			logger.Add(vlog.Fatal, types.ErrProxy, fmt.Errorf("can't start proxy %w", err))
		}
	}()

	logger.Add(vlog.Info, types.ResultOK, "the balancer is running")

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)
	<-stopSignal
}

// initializeConfig is the function that initializes the configuration of the application.
func initializeConfig() *config.Config {
	runtime.GOMAXPROCS(runtime.NumCPU())

	configFile := os.Getenv("ConfigFile")

	if configFile == "." {
		configFile = "../../config"
	}

	cfg := config.New()

	if err := cfg.Load(configFile); err != nil {
		log.Fatalf("%v", err)
	}

	if resultCode := cfg.InitProxyPort(); resultCode != types.ResultOK {
		log.Fatalf("can't init proxy: %s", resultCode.ToStr())
	}

	return cfg
}

// createPeerListForBalancer is the function that creates a list of peers for the balancer.
func createPeerListForBalancer(cfg *config.Config) []peer.IPeer {
	listPeer := make([]peer.IPeer, len(cfg.Peers))

	for index, valPeer := range cfg.Peers {
		listPeer[index] = valPeer
	}

	return listPeer
}
