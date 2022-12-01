package proxy

type Config struct {
	DefaultPort     string `yaml:"defaultPort"`
	TimeDeadLineMS  uint   `yaml:"timeDeadLineMs"`
	ShutdownTimeout uint   `yaml:"shutdownTimeout"`
}
