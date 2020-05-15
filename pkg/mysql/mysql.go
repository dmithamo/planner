package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // mysql db drivers - indirect
)

// Mysql is the interface thru which packages using this package access its services
type Mysql struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	IDB         *sql.DB
}

// OpenDB establishes a conn to a mysql db
func (m *Mysql) OpenDB(dsn string) (*sql.DB, error) {

	if dsn == "" {
		err := errors.New("invalid db connection url")
		m.ErrorLogger.Fatalf("[db] %v", err)
		return nil, err
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		m.ErrorLogger.Fatalf("[db] %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		m.ErrorLogger.Fatalf("[db] %v", err)
		return nil, err
	}
	m.InfoLogger.Println("[db] succesfully connected to db")
	return db, err
}

// CreateTables creates tables in the db.
// References schemas.go for ... schemas
func (m *Mysql) CreateTables() {
	for tableName, schema := range schemas {
		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (%v)", tableName, schema)
		m.InfoLogger.Printf("[db] %v", query)

		_, err := m.IDB.Exec(query)
		if err != nil {
			m.ErrorLogger.Fatalf("[db] %v", err)
		}
		m.InfoLogger.Printf("[db] succesfully created table: %v", tableName)
	}
}

// DropTables drops all the tables.
// Reference schemas.go
func (m *Mysql) DropTables() {
	for tableName := range schemas {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %v", tableName)
		m.InfoLogger.Printf("[db] %v", query)

		_, err := m.IDB.Exec(query)
		if err != nil {
			m.ErrorLogger.Fatalf("[db] %v", err)
		}
		m.InfoLogger.Printf("[db] succesfully dropped table: %v", tableName)
	}
}
