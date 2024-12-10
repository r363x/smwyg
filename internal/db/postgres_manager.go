package db

import (
    "database/sql"
	_ "github.com/lib/pq"
	"github.com/r363x/dbmanager/internal/config"
    "fmt"
)

// PostgreSQLManager represents a manager for interacting with PostgreSQL databases.
type PostgreSQLManager struct {
    db     *sql.DB
    cfg    config.DatabaseConfig
}

func NewPostgreSQLManager(cfg config.DatabaseConfig) (*PostgreSQLManager,
error) {
    return &PostgreSQLManager{cfg: cfg}, nil
}

func (m *PostgreSQLManager) DbType() string {
    return m.cfg.Type
}

func (m *PostgreSQLManager) Status() error {
    if err := m.db.Ping(); err != nil {
        return fmt.Errorf("failed to ping PostgreSQL database: %v", err)
    }
    return nil
}

// Connect establishes a connection to the PostgreSQL database.
func (m *PostgreSQLManager) Connect() error {
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        m.cfg.Host, m.cfg.Port, m.cfg.User, m.cfg.Password, m.cfg.DBName,
    )
    var err error
    m.db, err = sql.Open("postgres", dsn)
    if err != nil {
        return fmt.Errorf("failed to connect to PostgreSQL database: %v", err)
    }
    // Check the connection by pinging it.
    if err := m.Status(); err != nil {
        m.Disconnect()
        return err
    }
    return nil
}

// Disconnect closes the current PostgreSQL database connection.
func (m *PostgreSQLManager) Disconnect() error {
    if m.db == nil {
        return nil // Already disconnected or never connected.
    }
    var err error
    err = m.db.Close()
    m.db = nil
    return err
}

// ExecuteQuery executes a SQL query against the PostgreSQL database and returns its results as maps of string to interface{} values.
func (m *PostgreSQLManager) ExecuteQuery(query string) ([]map[string]interface{}, error) {
    rows, err := m.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to execute query against PostgreSQL: %v", err)
    }
    defer func() { _ = rows.Close() }()

    columns, err := rows.Columns()
    if err != nil {
        return nil, fmt.Errorf("failed to get column names from result set for PostgreSQL query execution: %v", err)
    }

    var result []map[string]interface{}
    valuePointers := make([]interface{}, len(columns))
    for i := range columns {
        // Use a generic interface{} as the type of each pointer, so that sql.Scan can handle any data type.
        valuePointers[i] = new(interface{})
    }

    for rows.Next() {
        row := make(map[string]interface{})
        err = rows.Scan(valuePointers...)
        if err != nil {
            return nil, fmt.Errorf("failed to scan result set for PostgreSQL query execution: %v", err)
        }

        // Populate map with actual values from value pointers and column names.
        for i, col := range columns {
            row[col] = *(valuePointers[i].(*interface{}))
        }
        result = append(result, row)
    }

    return result, nil
}

// GetTables returns a list of tables available in the current PostgreSQL database.
func (m *PostgreSQLManager) GetTables() ([]string, error) {
    rows, err := m.db.Query("SELECT table_name FROM information_schema.tables WHERE table_type='BASE TABLE'")
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
            WHERE table_catalog = '%s' AND table_name = '%s'
        `, "public", tableName),
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

