package proxy

type Config struct {
	Addr              string `yaml:"Addr"`
	ShutdownTimeout   uint   `yaml:"ShutdownTimeout"`
	ReadHeaderTimeout uint   `yaml:"ReadHeaderTimeout"`
	WriteTimeout      uint   `yaml:"WriteTimeout"`
	ReadTimeout       uint   `yaml:"ReadTimeout"`
	IdleTimeout       uint   `yaml:"IdleTimeout"`
}

