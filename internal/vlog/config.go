package vlog

type Config struct {
	DirLog         string `yaml:"dirLog"`
	FileSize       uint64 `yaml:"fileSize"`
	APIShowRecords uint64 `yaml:"apiShowRecords"`
	LogWriteSec    uint64 `yaml:"logWriteSec"`
	KindType       string `yaml:"kindType"`
}
