package db

import (
    "fmt"
	"github.com/r363x/dbmanager/internal/config"

)

// PostgreSQLManager represents a manager for interacting with PostgreSQL databases.
type PostgreSQLManager struct {
    BaseDBManager
}

func NewPostgreSQLManager(cfg config.DatabaseConfig) (*PostgreSQLManager,
error) {
    return &PostgreSQLManager{BaseDBManager{cfg: cfg}}, nil
}

// Connect establishes a connection to the PostgreSQL database.
func (m *PostgreSQLManager) Connect() error {
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        m.cfg.Host, m.cfg.Port, m.cfg.User, m.cfg.Password, m.cfg.DBName,
    )
    return m.BaseDBManager.Connect(dsn)
}

// ExecuteQuery executes a SQL query against the PostgreSQL database and returns its results as maps of string to interface{} values.
func (m *PostgreSQLManager) ExecuteQuery(query string) ([]map[string]interface{}, error) {
    rows, err := m.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to execute query against PostgreSQL: %v", err)
    }
    defer rows.Close()

    columns, err := rows.Columns()
    if err != nil {
        return nil, fmt.Errorf("failed to get column names from result set for PostgreSQL query execution: %v", err)
    }

	var result []map[string]interface{}
	for rows.Next() {
		row := make(map[string]interface{})
		values := make([]interface{}, len(columns))
		valuePointers := make([]interface{}, len(columns))
		for i := range columns {
			valuePointers[i] = &values[i]
		}

		if err := rows.Scan(valuePointers...); err != nil {
			return nil, err
		}

		for i, col := range columns {
			row[col] = values[i]
		}
		result = append(result, row)
	}

    return result, nil
}

// GetTables returns a list of tables available in the current PostgreSQL database.
func (m *PostgreSQLManager) GetTables() ([]string, error) {
    rows, err := m.db.Query("SELECT tablename FROM pg_tables t WHERE t.tableowner = current_user")
    if err != nil {
        return nil, fmt.Errorf("failed to query for tables in PostgreSQL: %v", err)
    }
    defer func() { _ = rows.Close() }()

    var tables []string
    for rows.Next() {
        var tableName string
        if err := rows.Scan(&tableName); err != nil {
            return nil, fmt.Errorf("error scanning result set while fetching table names from PostgreSQL: %v", err)
        }
        tables = append(tables, tableName)
    }

    return tables, nil
}

// GetTableStructure fetches the structure of a given table in the current database.
func (m *PostgreSQLManager) GetTableStructure(tableName string, databaseName string) (*TableStructure, error) {
    rows, err := m.db.Query(
        fmt.Sprintf(`
            SELECT column_name, data_type FROM information_schema.columns
            WHERE table_schema = 'public' AND table_name = '%s'
        `, tableName),
    )
    if err != nil {
        return &TableStructure{}, fmt.Errorf("failed to fetch table structure for PostgreSQL: %v", err)
    }
    defer func() { _ = rows.Close() }()

    var structure TableStructure
    for rows.Next() {
        var columnName string
        var dataType string

        // Scan each row.
        if err := rows.Scan(&columnName, &dataType); err != nil {
            return &TableStructure{}, fmt.Errorf("failed to scan result set while fetching PostgreSQL table: %v", err)
        }

        switch dataType {
        case "integer":
            dataType = "int"
        case "character varying":
            dataType = "varchar"
        }

        structure.Columns = append(structure.Columns, Column{
            Name:     columnName,
            DataType: dataType,
        })
    }
    return &structure, nil
}

func (m *PostgreSQLManager) GetVersion() (string, error) {
    var ver string
    err := m.db.QueryRow("SHOW server_version").Scan(&ver)
    if err != nil {
        return "", fmt.Errorf("failed to get PostgreSQL version: %v", err)
    }
    return ver, nil
}

func (m *PostgreSQLManager) GetDatabases() ([]string, string, error) {
	rows, err := m.db.Query("SELECT datname FROM pg_database")
    if err != nil {
        return nil, "",  err
    }
    defer rows.Close()

    var databases []string
    for rows.Next() {
        var database string
        if err = rows.Scan(&database); err != nil {
            return nil, "", err
        }
        databases = append(databases, database)
    }


	return databases, m.cfg.DBName, nil
}
