package proxy

type Config struct {
	DefaultPort       string `yaml:"DefaultPort"`
	ShutdownTimeout   uint   `yaml:"ShutdownTimeout"`
	ReadHeaderTimeout uint   `yaml:"ReadHeaderTimeout"`
	WriteTimeout      uint   `yaml:"WriteTimeout"`
	ReadTimeout       uint   `yaml:"ReadTimeout"`
	IdleTimeout       uint   `yaml:"IdleTimeout"`
}
