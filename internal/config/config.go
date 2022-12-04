package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/proxy"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"

	"gopkg.in/yaml.v2"
)

const defaultProxyPort = 8080

type Config struct {
	Logger         *vlog.Config         `yaml:"logger"`
	Proxy          *proxy.Config        `yaml:"proxy"`
	Peers          []*peer.Peer         `yaml:"peers"`
	CheckTimeAlive *peer.CheckTimeAlive `yaml:"peerCheckTimeAlive"`
	ProxyPort      string
}

func New() *Config {
	cfg := &Config{
		Logger:         nil,
		Proxy:          nil,
		Peers:          nil,
		CheckTimeAlive: nil,
		ProxyPort:      "",
	}

	return cfg
}

func (c *Config) Init() types.ResultCode {
	c.ProxyPort = fmt.Sprintf(":%s", os.Getenv("ProxyPort"))
	if c.ProxyPort == ":" {
		c.ProxyPort = fmt.Sprintf(":%d", defaultProxyPort)
	}

	if c.ProxyPort == strings.Trim(":", " ") {
		return types.ErrEmptyValue
	}

	return types.ResultOK
}

func (c *Config) Open(configFileName string) error {
	searchPathConfig := []string{configFileName, "", "./config/"}

	var isPathFound bool

	for _, searchPath := range searchPathConfig {
		configFilePath := filepath.Join(searchPath, "config.yaml")
		if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
			continue
		}

		isPathFound = true
		configFileName = configFilePath

		break
	}

	if !isPathFound {
		//nolint:goerr113
		return fmt.Errorf("failed: %s", types.ErrCantFindFile.ToStr())
	}

	fileConfig, err := os.Open(configFileName)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Fatalf("Error can't close config file: %s, err: %s", configFileName, err)
		}
	}(fileConfig)

	err = c.decodeConfigFileYaml(fileConfig)
	if err != nil {
		return fmt.Errorf("can't decode config file: %s, err: %w", configFileName, err)
	}

	return nil
}

//nolint
func (c *Config) decodeConfigFileYaml(configYaml *os.File) error {
	decoder := yaml.NewDecoder(configYaml)
	err := decoder.Decode(c)

	if err != nil {
		return fmt.Errorf("failed to decode config yml file: %w", err)
	}

	return nil
}
