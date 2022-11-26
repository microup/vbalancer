package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"vbalancer/internal/config"
	"vbalancer/internal/proxy"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

func main() {
	proxyPort := fmt.Sprintf(":%s", os.Getenv("ProxyPort"))
	if proxyPort == ":" {
		log.Fatalf("Can't read environment variable ProxyPort")
	}

	configFile :=os.Getenv("ConfigFile")
	if configFile == "" {
		log.Fatalf("Can't read environment variable ConfigFile")
	}

	ctx, serverCancel := context.WithCancel(context.Background())

	cfg, err := config.New(configFile)
	if err != nil {
		log.Fatalf("Can't create and init config from file: %s, err: %v", configFile, err)
	}

	logger, err := vlog.New(cfg.Logger)
	if err != nil {
		log.Fatalf("Error create logger: %v", err)
	}
	defer logger.Close()

	logger.Add(vlog.Info, types.ResultOK, "the balancer was running")
	proxy := proxy.New(ctx, proxyPort, cfg.Proxy, cfg.Peers, logger)

	go func() {
		logger.Add(vlog.Info, types.ResultOK, fmt.Sprintf("start server addr on %s", proxyPort))
		if err := proxy.Start(cfg.CheckTimeAlive); err != nil {
			logger.Add(vlog.Fatal, types.ResultOK, fmt.Sprintf("can't start server %s", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	serverCancel()

	var wg sync.WaitGroup
	wg.Add(1)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Proxy.ShutdownTimeout)*time.Second)
	defer cancel()

	go func() {
		defer wg.Done()
		if err := proxy.Shutdown(shutdownCtx); err != nil {
			logger.Add(vlog.Fatal, types.ResultOK, fmt.Sprintf("can't stop server %s", err))
		}
	}()

	wg.Wait()
}
