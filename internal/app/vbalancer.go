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
func Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx := context.Background()

	cfg := config.New()
	if err := cfg.Init(); err != nil {
		log.Panicf("failed to initialize configuration: %s", err.Error())
	}

	logger := vlog.New(cfg.Log)
	if err := logger.Init(); err != nil {
		log.Panicf("failed to create logger: %s", err)
	}
	defer logger.Close()

	defer func() {
		if err := recover(); err != nil {
			msgErr := fmt.Errorf("%w: %v", ErrRecoveredPanic, err)

			logger.Add(types.Fatal, types.ErrRecoverPanic, msgErr)
			log.Printf("%v", msgErr)
		}
	}()

	proxy, err := core.YamlToObject(cfg.Proxy, proxy.New())
	if err != nil {
		logger.Add(types.Fatal, types.ErrCantGetProxyObject, "can't get proxy object")
		log.Panicf("can't get proxy object, code: %d", types.ErrCantGetProxyObject)
	}

	if err = proxy.Init(ctx, logger); err != nil {
		logger.Add(types.Fatal, types.ErrCantInitProxy, fmt.Errorf("%w", err))
		log.Panicf("%v, code: %d", err, types.ErrCantInitProxy)
	}

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

	chanListenProxy := make(chan error)

	go func() {
		logger.Add(types.Info, types.ResultOK, fmt.Sprintf("start server addr on %s", proxy.Port))
		chanListenProxy <- proxy.ListenAndServe(ctx)

		stopSignal <- syscall.SIGTERM
	}()

	select {
	case <-ctx.Done():
		logger.Add(types.Info, types.ResultOK, "get ctx.Done()...")

	case listenErr := <-chanListenProxy:
		{
			if listenErr != nil {
				logger.Add(types.Fatal, types.ErrProxy, fmt.Errorf("the proxy was return err: %w", err))
			} else {
				log.Printf("the proxy was close")
			}
		}
	case <-stopSignal:
		{
			logger.Add(types.Info, types.ResultOK, "syscall.SIGTERM...")
		}
	}
}
