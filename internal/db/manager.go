package db

import (
	"github.com/r363x/dbmanager/internal/config"
    "fmt"
)

type Manager interface {
	Connect() error
	Status() error
    GetVersion() (string, error)
    GetDatabases() ([]string, string, error)
	Disconnect() error
	ExecuteQuery(query string) ([]map[string]interface{}, error)
	GetTables() ([]string, error)
	GetColumns(table string) ([]string, error)
}

func NewManager(cfg config.DatabaseConfig) (Manager, error) {
	switch cfg.Type {
	case "mysql":
		return NewMySQLManager(cfg)
	// case "postgres":
	//     return NewPostgresManager(cfg)
	// Add more database types
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}
}
