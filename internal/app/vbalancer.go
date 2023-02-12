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

	"vbalancer/internal/config"
	"vbalancer/internal/proxy"
	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

// Run this is the function of an application that starts a proxy server.
func Run(wgStartApp *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("catch err: %v", err) //nolint:forbidigo
		}
	}()

	// The function first initializes the configuration by calling `initializeConfig()`	
	configuration := initializeConfig()

	// and then creates a logger instance using `vlog.New(configuration.Logger)`
	logger, err := vlog.New(configuration.Logger)
	if err != nil {
		log.Panicf("failed created logger: %v", err)
	}

	defer func(logger vlog.ILog) {
		err = logger.Close()
		if err != nil {
			log.Fatalf("failed closed logger: %v", err)
		}
	}(logger)

	// The function creates a context with a cancel function for proxy work 
	ctx, proxyWorkCancel := context.WithCancel(context.Background())
	// and another context with a timeout for shutting down the application
	_, cancel := context.WithTimeout(context.Background(), configuration.Proxy.ShutdownTimeout)

	defer cancel()

	// The function then creates a list of peers for the proxy balancer using `createPeerListForBalancer(configuration)`
	peerList := createPeerListForBalancer(configuration)

	proxyBalancer := proxy.New(configuration.Proxy, peerList, logger)

	logger.Add(types.Info, types.ResultOK, fmt.Sprintf("start server addr on %s", configuration.ProxyPort))

	// The proxy server is started by calling `proxyBalancer.ListenAndServe(ctx, configuration.ProxyPort)`
	// in a separate goroutine
	go func() {
		if err = proxyBalancer.ListenAndServe(ctx, configuration.ProxyPort); err != nil {
			logger.Add(types.Fatal, types.ErrProxy, fmt.Errorf("can't start proxy %w", err))
		}
	}()

	logger.Add(types.Info, types.ResultOK, "the balancer is running")

	if wgStartApp != nil {
		wgStartApp.Done()
	}

	stopSignal := make(chan os.Signal, 1)
	// The function waits for the proxy balancer to be stopped by calling `proxyBalancer.Close()`
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)
	// and then waits for the proxy balancer to be stopped by calling `<-stopSignal`
	<-stopSignal

	// The function closes the proxy balancer
	proxyWorkCancel()
}

// initializeConfig initializes the configuration for the vbalancer application.
// It sets the maximum number of processors to be used by the Go runtime to the number of CPUs.
// It loads the configuration from the file specified in the ConfigFile environment variable, 
// or from the "../../config" file if the environment variable is not set.
// It initializes the proxy port and returns the configuration if successful, otherwise it logs an error and exits.
func initializeConfig() *config.Config {
	// Set the maximum number of processors to be used by the Go runtime to the number of CPUs.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Get the configuration file path from the ConfigFile environment variable.
	configFile := os.Getenv("ConfigFile")

	// If the environment variable is not set, use the default file "../../config".
	if configFile == "." {
		configFile = "../../config"
	}

	// Create a new configuration instance.	
	cfg := config.New()

	// Load the configuration from the file.	
	if err := cfg.Load(configFile); err != nil {
		// If there is an error loading the configuration, log the error and exit.
		log.Fatalf("%v", err)
	}

	// Initialize the proxy port.
	if resultCode := cfg.InitProxyPort(); resultCode != types.ResultOK {
		// If there is an error initializing the proxy port, log the error and exit.		
		log.Fatalf("can't init proxy: %s", resultCode.ToStr())
	}

	// Return the configuration.
	return cfg
}

// createPeerListForBalancer takes a Config struct and returns a slice of IPeer interface
// The function iterates over the list of Peers stored in the Config struct and converts
// each individual Peer into an IPeer interface.
// The resulting slice of IPeer interfaces is returned.
func createPeerListForBalancer(cfg *config.Config) []peer.IPeer {
	listPeer := make([]peer.IPeer, len(cfg.Peers))

	// Iterate over the Peers stored in the Config struct	
	for index, valPeer := range cfg.Peers {
		// Convert the individual Peer into an IPeer interface
		listPeer[index] = valPeer
	}

	// Return the resulting slice of IPeer interfaces
	return listPeer
}
