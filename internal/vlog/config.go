package vlog

type Config struct {
	DirLog            string `yaml:"DirLog"`
	FileSize          uint64 `yaml:"FileSize"`
	ApiShowRecords    uint64 `yaml:"ApiShowRecords"`	 
	LogWriteSec       uint64 `yaml:"LogWriteSec"`
	KindType          string `yaml:"KindType"`
}
