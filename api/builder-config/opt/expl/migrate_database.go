package main

import (
	"database/sql"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	"log"
)

//go:embed migrations/*.sql
var fs embed.FS

func MigrateDatabase(db *sql.DB) error {
	drv, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	defer drv.Close() // errors while closing are ignored

	src, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}
	defer src.Close() // errors while closing are ignored

	mig, err := migrate.NewWithInstance("iofs", src, "postgres", drv)
	if err != nil {
		return err
	}
	mig.Log = newLoggerAdapter()

	err = mig.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func newLoggerAdapter() migrate.Logger {
	return &loggerAdapter{}
}

type loggerAdapter struct{}

func (r *loggerAdapter) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (r *loggerAdapter) Verbose() bool {
	return true
}
