package database

import (
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

type Tx struct {
	tx *sqlx.Tx
}

// Multi Select
func (t *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := t.tx.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	return rows, nil
}

// Single Select
func (t *Tx) QueryRow(query string, args ...interface{}) *sql.Row {
	row := t.tx.QueryRow(query, args...)
	return row
}

func (t *Tx) Exec(query string, args ...interface{}) error {
	_, err := t.tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tx) NameExec(query string, args interface{}) error {
	_, err := t.tx.Exec(query, args)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tx) Commit() error {
	return t.tx.Commit()
}

func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

func (t *Tx) Get(record any, query string, args ...any) error {
	err := t.tx.Get(record, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tx) Select(record any, query string, args ...any) error {
	err := t.tx.Select(record, query, args...)
	if err != nil {
		return err
	}
	return nil
}
