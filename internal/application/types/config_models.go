package types

// ConfigModel model
type ConfigModel struct {
	Secrets  []string `json:"secrets,omitempty"`
	Database EngineDB `json:"db,omitempty"`
}

type DBConfig struct {
	Dialect        string `json:"dialect"`
	Database       string `json:"database"`
	Host           string `json:"host"`
	Port           int    `json:"port"`
	User           string `json:"user"`
	Password       string `json:"password"`
	MaxConnections int    `json:"max_conns"`
}

// EngineDB model
type EngineDB struct {
	Writer *DBConfig `json:"writer,omitempty"`
	Reader *DBConfig `json:"reader,omitempty"`
}
