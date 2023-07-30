package config

import (
	_ "embed"

	"github.com/choyri/gugu/internal/configx"
	"github.com/choyri/gugu/internal/databasex"
	"github.com/choyri/gugu/internal/logx"
)

type Config struct {
	ListenAddr string

	LogX      logx.Config
	DatabaseX databasex.Config
}

var (
	//go:embed config.default.toml
	defaultConfigData []byte
)

func Init() (*Config, error) {
	var conf Config

	err := configx.Init(
		&conf,
		configx.WithDefaultData(defaultConfigData),
	)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
