package config

import (
	"errors"
	"fmt"
	"io"
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
const DefaultConfigFile = "config.yaml"

var ErrCantGetProxyPort = errors.New("can't get proxy port")

// Config is the configuration of the proxy server.
type Config struct {
	// Logger is the configuration for the logger.
	Logger *vlog.Config `yaml:"logger"`
	// Proxy is the configuration for the proxy server.
	Proxy *proxy.Config `yaml:"proxy"`
	// Peers is a list of peer configurations.
	Peers []peer.Peer `yaml:"peers"`
	// Rules is the configuration for the rules to proxy.
	Rules *rules.Rules `yaml:"rules"`
	// ProxyPort is the port for the proxy server.
	ProxyPort string
}

// New creates a new configuration for the vbalancer application.
func New() *Config {
	return &Config{
		Logger:    nil,
		Proxy:     nil,
		Peers:     nil,
		Rules:     nil,
		ProxyPort: "",
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

	if resultCode := c.GetProxyPortConfig(); resultCode != types.ResultOK {
		return fmt.Errorf("%w: %s", ErrCantGetProxyPort, resultCode.ToStr())
	}

	return nil
}

// GetProxyPortConfig get the proxy port to serverconfiguration.
func (c *Config) GetProxyPortConfig() types.ResultCode {
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
