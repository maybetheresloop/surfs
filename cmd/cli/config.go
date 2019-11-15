package main

import (
	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
)

type blockConfig struct {
	Host string
	Port uint
}

type metadataConfig struct {
	Host string
	Port uint
}

type config struct {
	BlockConf    blockConfig    `toml:"block-store"`
	MetadataConf metadataConfig `toml:"metadata-store"`
}

func getConfig(c *cli.Context) (*config, error) {
	path := c.String("config")
	var conf config

	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
