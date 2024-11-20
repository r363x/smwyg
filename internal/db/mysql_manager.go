package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/r363x/dbmanager/internal/config"
)

type MySQLManager struct {
	db  *sql.DB
	cfg config.DatabaseConfig
}

func NewMySQLManager(cfg config.DatabaseConfig) (*MySQLManager, error) {
	return &MySQLManager{cfg: cfg}, nil
}

func (m *MySQLManager) Connect() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", m.cfg.User, m.cfg.Password, m.cfg.Host, m.cfg.Port, m.cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	m.db = db
	return db.Ping()
}

func (m *MySQLManager) Disconnect() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

func (m *MySQLManager) ExecuteQuery(query string) ([]map[string]interface{}, error) {
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
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

func (m *MySQLManager) GetTables() ([]string, error) {
	rows, err := m.db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func (m *MySQLManager) GetColumns(table string) ([]string, error) {
	rows, err := m.db.Query(fmt.Sprintf("SHOW COLUMNS FROM %s", table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var column string
		var rest interface{}
		if err := rows.Scan(&column, &rest, &rest, &rest, &rest, &rest); err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	return columns, nil
}
