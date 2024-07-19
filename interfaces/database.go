package interfaces

import (
	"database/sql"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type SQLExecuter interface {
	Exec(query string, args ...interface{}) error
	Get(record any, query string, args ...any) error
	Select(record any, query string, args ...any) error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	NameExec(query string, args interface{}) error
}

type Tx interface {
	Commit() error
	Rollback() error
}
