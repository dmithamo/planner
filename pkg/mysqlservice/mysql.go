package mysqlservice

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // mysql db drivers - indirect
)

// These loggers log very request/response to/fro the db service
var infoLogger *log.Logger
var errorLogger *log.Logger

// OpenDB establishes a conn to a db
func OpenDB(dsn string, iLogger *log.Logger, eLogger *log.Logger) (*sql.DB, error) {
	//Attach loggers passed from parent package.
	// Used sparingly though, to maintain loosest coupling possible
	infoLogger = iLogger
	errorLogger = eLogger

	if dsn == "" {
		err := errors.New("invalid db connection url")
		errorLogger.Fatalf("[db] %v", err)
		return nil, err
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		errorLogger.Fatalf("[db] %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		errorLogger.Fatalf("[db] %v", err)
		return nil, err
	}

	return db, err
}

// CreateTables creates a table in the db.
// References schemas.go for ... schemas
func CreateTables(iDB *sql.DB) {
	for tableName, schema := range schemas {
		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (%v)", tableName, schema)
		infoLogger.Printf("[db] %v", query)

		_, err := iDB.Exec(query)
		if err != nil {
			errorLogger.Fatalf("[db] %v", err)
		}
		infoLogger.Printf("[db] succesfully created table: %v", tableName)
	}
}

// DropTables drops all the tables.
// Reference schemas.go
func DropTables(iDB *sql.DB) {
	for tableName := range schemas {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %v", tableName)
		infoLogger.Printf("[db] %v", query)

		_, err := iDB.Exec(query)
		if err != nil {
			errorLogger.Fatalf("[db] %v", err)
		}
		infoLogger.Printf("[db] succesfully dropped table: %v", tableName)
	}
}
