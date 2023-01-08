package main

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

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cfg := initConfig()

	logger, err := vlog.New(cfg.Logger)
	if err != nil {
		log.Fatalf("error create logger: %v", err)
	}

	defer func(logger *vlog.VLog) {
		err = logger.Close()
		if err != nil {
			log.Fatalf("error close logger: %v", err)
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

	proxyBalancer := proxy.New(cfg.ProxyPort, cfg.Proxy, listPeer, logger)

	go func() {
		logger.Add(vlog.Info, types.ResultOK, fmt.Sprintf("start server addr on %s", cfg.ProxyPort))

		if err = proxyBalancer.Start(ctx, cfg.CheckTimeAlive); err != nil {
			logger.Add(vlog.Fatal, types.ResultOK, fmt.Sprintf("can't start server %s", err))
		}
	}()

	logger.Add(vlog.Info, types.ResultOK, "the balancer is running")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	proxyWorkCancel()

	var syncSrvShutdown sync.WaitGroup

	syncSrvShutdown.Add(1)

	go func() {
		defer syncSrvShutdown.Done()

		if err = proxyBalancer.Shutdown(); err != nil {
			logger.Add(vlog.Fatal, types.ResultOK, fmt.Sprintf("can't stop server %s", err))
		}
	}()

	syncSrvShutdown.Wait()
}

func initConfig() *config.Config {
	configFile := os.Getenv("ConfigFile")

	if configFile == "." {
		configFile = "../../config"
	}

	cfg := config.New()

	if resultCode := cfg.Init(); resultCode != types.ResultOK {
		log.Fatalf("can't init config: %s, err: %s", configFile, resultCode.ToStr())
	}

	if err := cfg.Load(configFile); err != nil {
		log.Fatalf("%v", err)
	}

	return cfg
}
