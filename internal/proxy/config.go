package proxy

type Config struct {
	DefaultPort            string `yaml:"defaultPort"`
	DeadLineTimeMS         uint   `yaml:"deadLineTimeMs"`
	ShutdownTimeoutSeconds uint   `yaml:"shutdownTimeoutSeconds"`
	SizeCopyBufferIO       uint   `yaml:"sizeCopyBufferIo"`
}
