package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
	"vbalancer/internal/config"
	"vbalancer/internal/proxy"
	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

func Run(wgStartApp *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("catch err: %v", err) //nolint:forbidigo
		}
	}()

	cfg := initConfig()

	logger, err := vlog.New(cfg.Logger)
	if err != nil {
		log.Panicf("failed created logger: %v", err)
	}

	defer func(logger vlog.ILog) {
		err = logger.Close()
		if err != nil {
			log.Fatalf("failed closed logger: %v", err)
		}
	}(logger)

	ctx, proxyWorkCancel := context.WithCancel(context.Background())
	_, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Proxy.ShutdownTimeoutSeconds)*time.Second)

	defer cancel()

	listPeer := createPeerList(cfg)

	proxyBalancer := proxy.New(cfg.Proxy, listPeer, logger)

	logger.Add(types.Info, types.ResultOK, fmt.Sprintf("start server addr on %s", cfg.ProxyPort))

	go func() {
		if err = proxyBalancer.Start(ctx, cfg.ProxyPort, cfg.CheckTimeAlive); err != nil {
			logger.Add(types.Fatal, types.ErrProxy, fmt.Sprintf("can't start proxy %s", err))
		}
	}()

	logger.Add(types.Info, types.ResultOK, "the balancer is running")

	if wgStartApp != nil {
		wgStartApp.Done()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	proxyWorkCancel()
}

func initConfig() *config.Config {
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

func createPeerList(cfg *config.Config) []peer.IPeer {
	listPeer := make([]peer.IPeer, len(cfg.Peers))
	
	for index, valPeer := range cfg.Peers {
		valPeer.Mu = &sync.RWMutex{}
		listPeer[index] = valPeer 
	}

	return listPeer
}