package proxy

type Config struct {
	DefaultPort            string `yaml:"defaultPort"`
	DeadLineTimeSeconds    uint   `yaml:"deadLineTimeSeconds"`
	ShutdownTimeoutSeconds uint   `yaml:"shutdownTimeoutSeconds"`
	SizeCopyBufferIO       uint   `yaml:"sizeCopyBufferIo"`
	ConnectionSemaphore    uint   `yaml:"connectionSemaphore"`
}
