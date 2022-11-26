package config

import (
	"errors"
	"log"
	"os"

	"github.com/microup/vbalancer/internal/peer"
	"github.com/microup/vbalancer/internal/proxy"
	"github.com/microup/vbalancer/internal/vlog"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Logger *vlog.Config  `yaml:"Logger"`
	Proxy  *proxy.Config `yaml:"Proxy"`
	Peers  []*peer.Peer  `yaml:"Peers"`
	CheckTimeAlive *peer.CheckTimeAlive `yaml:"PeerCheckTimeAlive"`
}

func New(configFileName string) (*Config, error) {
	if _, err := os.Stat(configFileName); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	f, err := os.Open(configFileName)
	if err != nil {
		return nil, err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("Error can't close config file: %s, err: %s", configFileName, err)
		}
	}(f)

	cfg := &Config{}

	err = cfg.decodeYaml(f)
	if err != nil {
		log.Fatalf("Can't decode config file: %s, err: %s", configFileName, err)
		return nil, err
	}

	return cfg, nil
}

func (c *Config) decodeYaml(configYaml *os.File) error {
	decoder := yaml.NewDecoder(configYaml)
	err := decoder.Decode(c)
	if err != nil {
		return err
	}
	return nil
}
