package proxy

type Config struct {
	DefaultPort                string `yaml:"defaultPort"`
	ClientDeadLineTimeSec      uint   `yaml:"clientDeadLineTimeSec"`
	DestinationHostTimeOutMs   uint   `yaml:"destinationHostTimeoutMs"`
	DestinationHostDeadLineSec uint   `yaml:"destinationHostDeadLineSec"`
	ShutdownTimeoutSeconds     uint   `yaml:"shutdownTimeoutSeconds"`
	ConnectionSemaphore        uint   `yaml:"connectionSemaphore"`
}
