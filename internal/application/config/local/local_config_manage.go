package local

import (
	"mossT8.github.com/device-backend/internal/application/config"
	"mossT8.github.com/device-backend/internal/application/config/types"
)

type Config struct {
}

func NewLocalConfigManager() config.Config {
	return &Config{}
}

func (sh *Config) GetConfig(_ string) (*types.ConfigModel, error) {
	db := &types.DBConfig{
		Dialect:        "mysql",
		Database:       "dishRatingsDB",
		Host:           "localhost",
		Port:           3306,
		User:           "root",
		Password:       "secret",
		MaxConnections: 150,
	}
	basePath := "/"
	return &types.ConfigModel{
		Server: types.HTTPServerConfig{
			ServerPort: 3001,
			BasePath:   &basePath,
		},

		Database: types.EngineDB{
			Writer: db,
			Reader: db,
		},
		BasePath: basePath,
	}, nil
}
