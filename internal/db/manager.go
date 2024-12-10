package db

import (
	"github.com/r363x/dbmanager/internal/config"
    "fmt"
    "strings"
)

type Manager interface {
	Connect() error
	DbType() string
    DbAddr() string
    DbUser() string
	Status() error
    GetVersion() (string, error)
    GetDatabases() ([]string, string, error)
	Disconnect() error
	ExecuteQuery(query string) ([]map[string]interface{}, error)
	GetTables() ([]string, error)
    GetTableStructure(tableName string, database string) (*TableStructure, error)
}

func NewManager(cfg config.DatabaseConfig) (Manager, error) {
	switch strings.ToLower(cfg.Type) {
	case "mysql","mariadb":
		return NewMySQLManager(cfg)
    case "postgresql":
        return NewPostgreSQLManager(cfg)
	// Add more database types
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}
}
