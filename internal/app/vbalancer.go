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
			msgErr := fmt.Errorf("%w: %v", types.ErrRecoveredPanic, err)

			logger.Add(vlog.Fatal, types.ErrGotPanic, msgErr)
			log.Printf("%v", msgErr)
		}
	}()

	defer func(logger vlog.ILog) {
		err = logger.Close()
		if err != nil {
			log.Fatalf("failed close logger: %v", err)
		}
	}(logger)

	proxyBalancer := cfg.Proxy

	err = proxyBalancer.Init(ctx, logger)
	if err != nil {
		logger.Add(vlog.Fatal, types.ErrCantInitProxy, fmt.Errorf("%w: %v", types.ErrInitProxy, err))
	}

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

	listenProxyChan := make(chan error)

	go func() {
		logger.Add(vlog.Info, types.ResultOK, fmt.Sprintf("start server addr on %s", cfg.Proxy.Port))
		listenProxyChan <- proxyBalancer.ListenAndServe(ctx, cfg.Proxy.Port)

		stopSignal <- syscall.SIGTERM
	}()

	listenErr := <-listenProxyChan
	if listenErr != nil {
		logger.Add(vlog.Fatal, types.ErrProxy, fmt.Errorf("the proxy was return err: %w", err))
	}

	<-stopSignal

	logger.Add(vlog.Info, types.ResultOK, "received shutdown signal, exiting gracefully...")
}
