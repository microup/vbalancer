package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"vbalancer/internal/config"
	"vbalancer/internal/core"
	"vbalancer/internal/proxy"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

var ErrRecoveredPanic = errors.New("recovered from panic")

// Run this is the function of an application that starts a proxy server.
func Run(ctx context.Context) {
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
		}
	}()

	proxy, err := core.YamlToObject(cfg.Proxy, proxy.New())
	if err != nil {
		logger.Add(types.Fatal, types.ErrCantGetProxyObject, "can't get proxy object")
	}

	if err = proxy.Init(ctx, logger); err != nil {
		logger.Add(types.Fatal, types.ErrCantInitProxy, fmt.Errorf("%w", err))
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	proxy.ListenAndServe(ctx)

	select {
	case s := <-interrupt:
		logger.Add(types.Info, types.ResultOK, "syscall.SIGTERM", s.String())
	case <-ctx.Done():
		logger.Add(types.Info, types.ResultOK, "Done.")
	case err = <-proxy.Notify():
		logger.Add(types.Fatal, types.ErrProxy, fmt.Errorf("the proxy returned err: %w", err))
	}
}
