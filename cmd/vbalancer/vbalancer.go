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
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cfg := initConfig()

	logger, err := vlog.New(cfg.Logger)
	if err != nil {
		log.Fatalf("Error create logger: %v", err)
	}
	defer logger.Close()

	ctx, serverCancel := context.WithCancel(context.Background())
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Proxy.ShutdownTimeout)*time.Second)

	defer cancel()

	logger.Add(vlog.Info, types.ResultOK, "the balancer was running")
	proxy := proxy.New(cfg.ProxyPort, cfg.Proxy, cfg.Peers, logger)

	go func() {
		logger.Add(vlog.Info, types.ResultOK, fmt.Sprintf("start server addr on %s", cfg.ProxyPort))

		if err = proxy.Start(ctx, cfg.CheckTimeAlive); err != nil {
			logger.Add(vlog.Fatal, types.ResultOK, fmt.Sprintf("can't start server %s", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	serverCancel()

	var syncSrvShutdown sync.WaitGroup

	syncSrvShutdown.Add(1)

	go func() {
		defer syncSrvShutdown.Done()

		if err = proxy.Shutdown(shutdownCtx); err != nil {
			logger.Add(vlog.Fatal, types.ResultOK, fmt.Sprintf("can't stop server %s", err))
		}
	}()

	syncSrvShutdown.Wait()
}

func initConfig() *config.Config {
	configFile := os.Getenv("ConfigFile")

	if configFile == "" {
		log.Fatalf("Can't read environment variable ConfigFile")
	}

	cfg := config.New()

	resultCode := cfg.Init()
	if resultCode != types.ResultOK {
		log.Fatalf("can't init config: %s, err: %s", configFile, resultCode.ToStr())
	}

	err := cfg.Open(configFile)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return cfg
}
