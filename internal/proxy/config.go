package proxy

import "time"

type Config struct {
	DefaultPort             string        `yaml:"defaultPort"`
	ClientDeadLineTime      time.Duration `yaml:"clientDeadLineTime"`
	DestinationHostTimeOut  time.Duration `yaml:"destinationHostTimeout"`
	DestinationHostDeadLine time.Duration `yaml:"destinationHostDeadLine"`
	ShutdownTimeout         time.Duration `yaml:"shutdownTimeout"`
	ConnectionSemaphore     uint          `yaml:"connectionSemaphore"`
	CountDialAttemptsToPeer uint          `yaml:"countDialAttemptsToPeer"`
}
