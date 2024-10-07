package config

import "mossT8.github.com/device-backend/internal/application/types"

type Config interface {
	GetConfig(configName string) (*types.ConfigModel, error)
}
