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
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"

	"gopkg.in/yaml.v2"
)

const DefaultProxyPort = 8080

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
	osEnvValue := os.Getenv("ProxyPort")
	if osEnvValue == ":" {
		return types.ErrEmptyValue
	}

	c.ProxyPort = fmt.Sprintf(":%s", osEnvValue)
	if c.ProxyPort == ":" {
		c.ProxyPort = fmt.Sprintf(":%d", DefaultProxyPort)
	}

 	c.ProxyPort  =  strings.Trim(c.ProxyPort, " ")

	if c.ProxyPort == strings.Trim(":", " ") {
		return types.ErrEmptyValue
	}

	return types.ResultOK
}

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

	err = c.decodeConfigFileYaml(fileConfig)
	if err != nil {
		return fmt.Errorf("can't decode config file: %s, err: %w", cfgFileName, err)
	}

	return nil
}

func (c *Config) decodeConfigFileYaml(configYaml *os.File) error {
	decoder := yaml.NewDecoder(configYaml)
	err := decoder.Decode(c)

	if err != nil {
		return fmt.Errorf("failed to decode config yml file: %w", err)
	}

	return nil
}
