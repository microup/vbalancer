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
	"vbalancer/internal/core"
	"vbalancer/internal/proxy"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

var ErrRecoveredPanic = errors.New("recovered from panic")

// Run this is the function of an application that starts a proxy server.
//
//nolint:funlen,cyclop
func Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx := context.Background()

	cfg := config.New()
	logger := vlog.New(cfg.Log)

	err := cfg.Init()
	if err != nil {
		log.Panicf("failed to initialize configuration: %s", err.Error())
	}

	err = logger.Init()
	if err != nil {
		log.Panicf("failed to create logger: %s", err)
	}

	defer func(logger vlog.ILog) {
		err = logger.Close()
		if err != nil {
			log.Fatalf("failed close logger: %v", err)
		}
	}(logger)

	defer func() {
		if err := recover(); err != nil {
			msgErr := fmt.Errorf("%w: %v", ErrRecoveredPanic, err)

			logger.Add(vlog.Fatal, types.ErrGotPanic, msgErr)
			log.Printf("%v", msgErr)
		}
	}()

	proxy, err := core.GetObjectFromMap(cfg.Proxy, proxy.New())
	if err != nil {
		logger.Add(vlog.Fatal, types.ErrCantGetProxyObject, "can't get proxy object")
	}

	err = proxy.Init(ctx, logger)
	if err != nil {
		logger.Add(vlog.Fatal, types.ErrCantInitProxy, fmt.Errorf("%w", err))
	}

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

	chanListenProxy := make(chan error)

	go func() {
		logger.Add(vlog.Info, types.ResultOK, fmt.Sprintf("start server addr on %s", proxy.Port))
		chanListenProxy <- proxy.ListenAndServe(ctx)

		stopSignal <- syscall.SIGTERM
	}()

	select {
	case <-ctx.Done():
		logger.Add(vlog.Info, types.ResultOK, "get ctx.Done()...")

	case listenErr := <-chanListenProxy:
		{
			if listenErr != nil {
				logger.Add(vlog.Fatal, types.ErrProxy, fmt.Errorf("the proxy was return err: %w", err))
			} else {
				log.Printf("the proxy was close")
			}
		}
	case <-stopSignal:
		{
			logger.Add(vlog.Info, types.ResultOK, "get syscall.SIGTERM...")
		}
	}
}
