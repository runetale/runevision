package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/runetale/runevision/domain/config"
	"github.com/runetale/runevision/utility"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db     *sqlx.DB
	url    string
	logger *utility.Logger
}

func NewPostgres(log *utility.Logger, url string) (*Postgres, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		db:     db,
		url:    url,
		logger: log,
	}, nil
}

func NewPostgresFromConfig(log *utility.Logger, cfg config.Postgres) (*Postgres, error) {
	return NewPostgres(log, fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DatabaseName))
}

func (d *Postgres) CreateDB(dbname string) error {
	_, err := d.db.Exec("create database " + dbname)
	// TODO: add error handling when already exist
	d.logger.Info(err.Error())
	return nil
}

func (d *Postgres) MigrateUp(databaseDir string) error {
	m, err := migrate.New(fmt.Sprintf("file://%s", databaseDir), d.url)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		d.logger.Error(fmt.Sprintf("migrate up error %s", err.Error()))
		return err
	}

	d.logger.Info("migrate up done with success")

	return nil
}

func (d *Postgres) MigrateDown(databaseDir string) error {
	m, err := migrate.New(fmt.Sprintf("file://%s", databaseDir), d.url)
	if err != nil {
		return err
	}

	if err = m.Down(); err != nil {
		d.logger.Error("migrate down error", err)
		return err
	}

	d.logger.Debug("migrrate down done with success")

	return err
}

func (d *Postgres) Ping() error {
	return d.db.Ping()
}

func (d *Postgres) Exec(query string, args ...interface{}) error {
	_, err := d.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (d *Postgres) NameExec(query string, args interface{}) error {
	_, err := d.db.NamedExec(query, args)
	if err != nil {
		return err
	}
	return nil
}

// Multi Select
func (d *Postgres) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	return rows, nil
}

// Single Select
func (d *Postgres) QueryRow(query string, args ...interface{}) *sql.Row {
	row := d.db.QueryRow(query, args...)
	return row
}

func (d *Postgres) Begin() (*Tx, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}

	return &Tx{tx: tx}, nil
}

func (d *Postgres) Get(record any, query string, args ...any) error {
	err := d.db.Get(record, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (d *Postgres) Select(record any, query string, args ...any) error {
	err := d.db.Select(record, query, args...)
	if err != nil {
		return err
	}
	return nil
}
