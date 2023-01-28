package proxy

type Config struct {
	DefaultPort      string `yaml:"defaultPort"`
	DeadLineTimeMS   uint   `yaml:"deadLineTimeMs"`
	ShutdownTimeout  uint   `yaml:"shutdownTimeout"`
	SizeCopyBufferIO uint   `yaml:"sizeCopyBufferIO"`
}
