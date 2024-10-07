package local

import (
	"mossT8.github.com/device-backend/internal/application/types"
	"mossT8.github.com/device-backend/internal/infrastructure/config"
)

type Config struct {
}

func NewLocalConfigManager() config.Config {
	return &Config{}
}

func (sh *Config) GetConfig(_ string) (*types.ConfigModel, error) {
	db := &types.DBConfig{
		Dialect:        "mysql",
		Database:       "deviceDB",
		Host:           "localhost",
		Port:           3306,
		User:           "root",
		Password:       "secret",
		MaxConnections: 150,
	}
	return &types.ConfigModel{
		Database: types.EngineDB{
			Writer: db,
			Reader: db,
		},
	}, nil
}
