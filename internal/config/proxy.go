package config

import "time"

// Config define a struct to hold the configuration.
type Proxy struct {
	// Define the default port to listen on
	DefaultPort string `yaml:"defaultPort"`
	// Define the client deadline time
	ClientDeadLineTime time.Duration `yaml:"clientDeadLineTime"`
	// Define the peer host timeout
	PeerHostTimeOut time.Duration `yaml:"peerHostTimeout"`
	// Define the peer host deadline
	PeerHostDeadLine time.Duration `yaml:"peerHostDeadLine"`
	// Define the max connection semaphore
	MaxCountConnection uint `yaml:"maxCountConnection"`
	// Define the count dial attempts to peer
	CountMaxDialAttemptsToPeer uint `yaml:"countDialAttemptsToPeer"`
}
