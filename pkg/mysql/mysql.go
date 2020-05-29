package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // mysql db drivers - indirect
)

// Mysql is the interface thru which packages using this package access its services
type Mysql struct {
	IDB *sql.DB
}

// OpenDB establishes a conn to a mysql db
func (m *Mysql) OpenDB(dsn string) (*sql.DB, error) {

	if dsn == "" {
		err := errors.New("invalid db connection url")
		return nil, err
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, err
}

// CreateTables creates tables in the db.
// References schemas.go for ... schemas
func (m *Mysql) CreateTables() error {
	for tableName, schema := range schemas {
		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v %v", tableName, schema)
		_, err := m.IDB.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

// DropTables drops all the tables.
// Reference schemas.go
func (m *Mysql) DropTables() error {
	for tableName := range schemas {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %v", tableName)
		_, err := m.IDB.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}
