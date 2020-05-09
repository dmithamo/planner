package mysqlservice

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql" // mysql db drivers - indirect
)

// iDB is an instance of the *sql.DB
var iDB *sql.DB

// Data exposes db data via methods defined on the models (tables) in the db
type Data struct {
	Projects Projects
}

// OpenDB establishes a conn to a db
func OpenDB(dsn string) (*sql.DB, error) {

	if dsn == "" {
		return nil, errors.New("invalid db connection url")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	iDB = db
	return db, err
}
