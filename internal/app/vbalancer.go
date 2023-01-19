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

//nolint:funlen
func Run(wgStartApp *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("catch err: %v", err) //nolint:forbidigo
		}
	}()

	cfg := initConfig()

	logger, err := vlog.New(cfg.Logger)
	if err != nil {
		log.Panicf("create logger: %v", err)
	}

	defer func(logger *vlog.VLog) {
		err = logger.Close()
		if err != nil {
			log.Fatalf("close logger: %v", err)
		}
	}(logger)

	ctx, proxyWorkCancel := context.WithCancel(context.Background())
	_, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Proxy.ShutdownTimeout)*time.Second)

	defer cancel()

	listPeer := make([]peer.IPeer, len(cfg.Peers))

	for i, v := range cfg.Peers {
		v.Mu = &sync.RWMutex{}
		listPeer[i] = v
	}

	proxyBalancer := proxy.New(cfg.Proxy, listPeer, logger)

	go func() {
		logger.Add(vlog.Info, types.ResultOK, fmt.Sprintf("start server addr on %s", cfg.ProxyPort))

		if err = proxyBalancer.Start(ctx, cfg.ProxyPort, cfg.CheckTimeAlive); err != nil {
			logger.Add(vlog.Fatal, types.ErrProxy, fmt.Sprintf("can't start proxy %s", err))
		}
	}()

	logger.Add(vlog.Info, types.ResultOK, "the balancer is running")

	if wgStartApp != nil {
		wgStartApp.Done()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	proxyWorkCancel()

	var srvShutdown sync.WaitGroup

	srvShutdown.Add(1)

	go func() {
		defer srvShutdown.Done()

		if err = proxyBalancer.Shutdown(); err != nil {
			logger.Add(vlog.Fatal, types.ResultOK, fmt.Sprintf("shutdown proxy err: %s", err))
		}
	}()

	srvShutdown.Wait()
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

	initProxy(cfg)

	return cfg
}

func initProxy(cfg *config.Config) {
	if resultCode := cfg.Init(); resultCode != types.ResultOK {
		log.Fatalf("can't init proxy: %s", resultCode.ToStr())
	}
}
