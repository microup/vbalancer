package proxy

type Config struct {
	DefaultPort         string `yaml:"defaultPort"`
	DeadLineTimeSeconds uint   `yaml:"deadLineTimeSeconds"`
	ShutdownTimeout     uint   `yaml:"shutdownTimeout"`
	SizeCopyBufferIO    uint   `yaml:"sizeCopyBufferIo"`
}
