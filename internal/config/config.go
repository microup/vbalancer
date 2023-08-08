package config

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"vbalancer/internal/proxy"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"

	"gopkg.in/yaml.v2"
)

const DefaultConfigFile = "config.yaml"
const DefaultFileLogSizeBytes = 100000
const DeafultShowRecordsAPI = 50
const DefaultDirLogs = "/logs"

// Config is the configuration of the proxy server.
type Config struct {
	// Logger is the configuration for the logger.
	Logger *vlog.Config `yaml:"logger" json:"logger"`
	// Proxy is the configuration for the proxy server.
	Proxy *proxy.Proxy `yaml:"proxy" json:"proxy"`
}

// New creates a new configuration for the vbalancer application.
func New() *Config {
	return &Config{
		Logger: &vlog.Config{
			DirLog:         DefaultDirLogs,
			FileSize:       DefaultFileLogSizeBytes,
			APIShowRecords: DeafultShowRecordsAPI,
		},
		Proxy: nil,
	}
}

func (c *Config) Init() error {
	configFile := os.Getenv("ConfigFile")
	if configFile == "" {
		configFile = DefaultConfigFile
	}

	if err := c.Load(configFile); err != nil {
		return err
	}

	if c.Proxy == nil {
		return fmt.Errorf("%w", types.ErrCantGetProxySection)
	}

	if resultCode := c.Proxy.UpdatePort(); resultCode != types.ResultOK {
		return fmt.Errorf("%w: %s", types.ErrCantGetProxyPort, resultCode.ToStr())
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
		//nolint:goerr113
		return fmt.Errorf("path to config not found: %s", cfgFileName)
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
