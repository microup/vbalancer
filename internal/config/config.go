package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"vbalancer/internal/proxy"
	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/proxy/rules"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"

	"gopkg.in/yaml.v2"
)

// DefaultProxyPort is the default port for the proxy server.
const DefaultProxyPort = 8080

// Configuration is the configuration for the vbalancer application.
type Config struct {
	// Logger is the configuration for the logger.
	Logger *vlog.Config `yaml:"logger"`
	// Proxy is the configuration for the proxy server.
	Proxy *proxy.Config `yaml:"proxy"`
	// Peers is a list of peer configurations.
	Peers []*peer.Peer `yaml:"peers"`
	// Rules is the configuration for the rules to proxy.
	Rules *rules.Rules `yaml:"rules"`
	// ProxyPort is the port for the proxy server.
	ProxyPort string
}

// New creates a new configuration for the vbalancer application.
func New() *Config {
	cfg := &Config{
		Logger:    nil,
		Proxy:     nil,
		Peers:     nil,
		Rules:     nil,
		ProxyPort: "",
	}

	return cfg
}

// InitializeConfig initializes the configuration for the vbalancer application.
func (c *Config) InitProxyPort() types.ResultCode {
	osEnvValue := os.Getenv("ProxyPort")
	if osEnvValue == ":" {
		return types.ErrEmptyValue
	}

	c.ProxyPort = fmt.Sprintf(":%s", osEnvValue)
	if c.ProxyPort == ":" {
		c.ProxyPort = fmt.Sprintf(":%d", DefaultProxyPort)
	}

	c.ProxyPort = strings.Trim(c.ProxyPort, " ")

	if c.ProxyPort == strings.Trim(":", " ") {
		return types.ErrEmptyValue
	}

	return types.ResultOK
}

// Load loads the configuration for the vbalancer application.
func (c *Config) Load(cfgFileName string) error {
	searchPathConfig := []string{cfgFileName, "", "./config/", "../../config/", "../config/", "../../../config"}

	var isPathFound bool

	for _, searchPath := range searchPathConfig {
		cfgFilePath := filepath.Join(searchPath, "config.yaml")
		if _, err := os.Stat(cfgFilePath); errors.Is(err, os.ErrNotExist) {
			continue
		}

		isPathFound = true
		cfgFileName = cfgFilePath

		break
	}

	if !isPathFound {
		//nolint:goerr113
		return fmt.Errorf("failed: %s", types.ErrCantFindFile.ToStr())
	}

	fileConfig, err := os.Open(cfgFileName)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Fatalf("Error can't close config file: %s, err: %s", cfgFileName, err)
		}
	}(fileConfig)

	err = c.DecodeConfigFileYaml(fileConfig)
	if err != nil {
		return fmt.Errorf("can't decode config file: %s, err: %w", cfgFileName, err)
	}

	return nil
}

// DecodeConfigFileYaml decodes the configuration for the vbalancer application.
func (c *Config) DecodeConfigFileYaml(configYaml *os.File) error {
	decoder := yaml.NewDecoder(configYaml)
	err := decoder.Decode(c)

	if err != nil {
		return fmt.Errorf("failed to decode config yml file: %w", err)
	}

	return nil
}
