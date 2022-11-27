package proxy

type Config struct {
	DefaultPort     string `yaml:"DefaultPort"`
	TimeDeadLineMS  uint   `yaml:"TimeDeadLineMS"`
	ShutdownTimeout uint   `yaml:"ShutdownTimeout"`
}
