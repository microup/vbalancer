package config

// Config defines the structure for storing configuration values read from a YAML file.
type Log struct {
	// directory where log files are stored
	DirLog string `yaml:"dirLog"`
	// maximum size of a log file, in megabytes
	FileSizeMB float64 `yaml:"fileSizeInMb"`
	// number of records to show when
	APIShowRecords uint64 `yaml:"apiShowRecords"`
}
