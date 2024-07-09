package config

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"vbalancer/internal/types"

	"gopkg.in/yaml.v2"
)

var ErrCantFindProxySection = errors.New("can't find proxy section in config")

// Config is the configuration of the proxy server.
type Config struct {
	// Logger is the configuration for the logger.
	Log *Log `json:"logger" yaml:"logger"`
	// Proxy is the configuration for the proxy server.
	Proxy any `json:"proxy" yaml:"proxy"`
}

// New creates a new configuration for the vbalancer application.
func New() *Config {
	return &Config{
		Log: &Log{
			DirLog:         types.DefaultDirLogs,
			FileSizeMB:     types.DefaultFileLogSizeMB,
			APIShowRecords: types.DeafultShowRecordsAPI,
		},
		Proxy: nil,
	}
}

// Init initializes the configuration by loading values from a YAML file.
func (c *Config) Init() error {
	configFile := os.Getenv("ConfigFile")
	if configFile == "" {
		configFile = types.DefaultNameConfigFile
	}

	if err := c.Load(configFile); err != nil {
		return err
	}

	if c.Proxy == nil {
		return fmt.Errorf("%w", ErrCantFindProxySection)
	}

	return nil
}

// Load loads the configuration for the vbalancer application.
func (c *Config) Load(cfgFileName string) error {
	searchPathConfig := []string{"", "./config/", "../../config/", "../config/", "../../../config"}

	var isPathFound bool

	for _, searchPath := range searchPathConfig {
		cfgFilePath := filepath.Join(searchPath, cfgFileName)

		info, err := os.Stat(cfgFilePath)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}

		if info.IsDir() {
			continue
		}

		isPathFound = true
		cfgFileName = cfgFilePath

		break
	}

	if !isPathFound {
		return fmt.Errorf("path to `%s` not found", cfgFileName)
	}

	fileConfig, err := os.Open(cfgFileName)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Fatalf("error can't close config file: %s, err: %s", cfgFileName, err)
		}
	}(fileConfig)

	err = c.decodeConfig(fileConfig)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// decodeConfig decodes the YAML configuration file.
func (c *Config) decodeConfig(configYaml io.Reader) error {
	decoder := yaml.NewDecoder(configYaml)

	err := decoder.Decode(c)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
