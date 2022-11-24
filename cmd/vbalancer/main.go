package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/microup/vbalancer/internal/config"
	"github.com/microup/vbalancer/internal/proxy"
	"github.com/microup/vbalancer/internal/types"
	"github.com/microup/vbalancer/internal/vlog"
)

func main() {
	cfgFile := flag.String("cfg_file", "config.yaml", "config file")
	flag.Parse()

	_, serverCancel := context.WithCancel(context.Background())

	cfg, err := config.New(*cfgFile)
	if err != nil {
		log.Fatalf("Can't create and init config from file: %s, err: %v", *cfgFile, err)
	}

	logger, err := vlog.New(cfg.Logger)
	if err != nil {
		log.Fatalf("Error create logger: %v", err)
	}
	defer logger.Close()

	logger.Add(vlog.Info, types.ResultOK, "the balancer was running")
	proxy := proxy.New(cfg.Proxy, logger, cfg.Peers)

	go func() {
		logger.Add(vlog.Info, types.ResultOK, fmt.Sprintf("start server addr on %s", cfg.Proxy.Addr))
		if err := proxy.Start(); err != nil {
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
