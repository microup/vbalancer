package vlog

// Config defines the structure for storing configuration values read from a YAML file.
type Config struct {
	 // directory where log files are stored
	DirLog         string `yaml:"dirLog"`
	// maximum size of a log file, in bytes
	FileSize       uint64 `yaml:"fileSize"`
	// number of records to show when 
	APIShowRecords uint64 `yaml:"apiShowRecords"`
}
