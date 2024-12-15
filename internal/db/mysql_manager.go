package db

import (
	"fmt"
	"github.com/r363x/dbmanager/internal/config"
)

type MySQLManager struct {
    BaseDBManager
}

type Column struct {
    Name string
    DataType string
}

type TableStructure struct {
    Columns []Column
}


func NewMySQLManager(cfg config.DatabaseConfig) (*MySQLManager, error) {
    return &MySQLManager{BaseDBManager{cfg: cfg}}, nil
}

func (m *MySQLManager) Connect() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
        m.cfg.User,
        m.cfg.Password,
        m.cfg.Host,
        m.cfg.Port,
        m.cfg.DBName,
    )
    return m.BaseDBManager.Connect(dsn)
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

func (m *MySQLManager) GetTableStructure(tableName string, database string) (*TableStructure, error) {
    rows, err := m.db.Query(
        fmt.Sprintf(`
            SELECT COLUMN_NAME, DATA_TYPE FROM information_schema.columns
            WHERE TABLE_SCHEMA = '%s' AND
            TABLE_NAME = '%s'`,
        database, tableName))
    if err != nil {
        return &TableStructure{}, err
    }

    defer rows.Close()

    var structure TableStructure

    for rows.Next() {
        var columnName string
        var dataType string

        if err := rows.Scan(&columnName, &dataType); err != nil {
            return &structure, err
        }

        structure.Columns = append(structure.Columns, Column{
            Name:  columnName,
            DataType: dataType,
        })
    }

    return &structure, nil
}

func (m *MySQLManager) GetVersion() (string, error) {
	row := m.db.QueryRow("SELECT VERSION()")

    var version string

    if err := row.Scan(&version); err != nil {
        return "", err
    }

	return version, nil
}


func (m *MySQLManager) GetDatabases() ([]string, string, error) {
	rows, err := m.db.Query("SHOW DATABASES")
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
