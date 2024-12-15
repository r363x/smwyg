package db

import (
	"database/sql"
	"fmt"
    "strings"
	"github.com/r363x/dbmanager/internal/config"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)


type BaseDBManager struct {
    db  *sql.DB
    cfg config.DatabaseConfig
}

func (b *BaseDBManager) Connect(dsn string) error {
    var err error
    b.db, err = sql.Open(strings.ToLower(b.cfg.Type), dsn)
    if err != nil {
        return fmt.Errorf("failed to connect to database: %v", err)
    }
    // Check the connection by pinging it.
    if err := b.Status(); err != nil {
        b.Disconnect()
        return err
    }
    return nil
}

func (b *BaseDBManager) Disconnect() error {
    if b.db == nil {
        return nil // Already disconnected or never connected.
    }
    var err error
    err = b.db.Close()
    b.db = nil
    return err
}

func (m *BaseDBManager) Status() error {
    return m.db.Ping()
}

func (m *BaseDBManager) DbType() string {
    return m.cfg.Type
}

func (m *BaseDBManager) DbAddr() string {
    return m.cfg.Host
}

func (m *BaseDBManager) DbUser() string {
    return m.cfg.User
}

