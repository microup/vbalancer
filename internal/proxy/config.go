package proxy

import "time"

// Config define a struct to hold the configuration.
type Config struct {
	// Define the default port to listen on
	DefaultPort             string        `yaml:"defaultPort"`
	// Define the client deadline time
	ClientDeadLineTime      time.Duration `yaml:"clientDeadLineTime"`
	// Define the destination host timeout
	DestinationHostTimeOut  time.Duration `yaml:"destinationHostTimeout"`
	// Define the destination host deadline
	DestinationHostDeadLine time.Duration `yaml:"destinationHostDeadLine"`
	// Define the shutdown timeout
	ShutdownTimeout         time.Duration `yaml:"shutdownTimeout"`
	// Define the connection semaphore
	ConnectionSemaphore     uint          `yaml:"connectionSemaphore"`
	// Define the count dial attempts to peer
	CountDialAttemptsToPeer uint          `yaml:"countDialAttemptsToPeer"`
}
