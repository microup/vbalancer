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
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

// Run this is the function of an application that starts a proxy server.
func Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
			logger.Add(vlog.Fatal, types.ErrGotPanic, fmt.Errorf("%w: %v", types.ErrRecoveredPanic, err))
		}
	}()

	defer func(logger vlog.ILog) {
		err = logger.Close()
		if err != nil {
			log.Fatalf("failed close logger: %v", err)
		}
	}(logger)

	proxyBalancer := proxy.New(cfg, cfg.Rules, logger)

	err = proxyBalancer.Init(cfg)
	if err != nil {
		logger.Add(vlog.Fatal, types.ErrCantInitProxy, fmt.Errorf("%w: %v", types.ErrInitProxy, err))
	}

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

	listenProxyChan := make(chan error)

	go func() {
		logger.Add(vlog.Info, types.ResultOK, fmt.Sprintf("start server addr on %s", cfg.ProxyPort))
		listenProxyChan <- proxyBalancer.ListenAndServe(ctx, cfg.ProxyPort)

		stopSignal <- syscall.SIGTERM
	}()

	listenErr := <-listenProxyChan
	if listenErr != nil {
		logger.Add(vlog.Fatal, types.ErrProxy, fmt.Errorf("the proxy was return err: %w", err))
	}

	<-stopSignal

	logger.Add(vlog.Info, types.ResultOK, "received shutdown signal, exiting gracefully...")
}
