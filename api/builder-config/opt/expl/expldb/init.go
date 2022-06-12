package expldb

import (
	"database/sql"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

func Init(databaseURL string) (*ExplDB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	waitUntilAvailable(db)
	err = applyMigrations(databaseURL)
	if err != nil {
		return nil, err
	}

	return &ExplDB{
		db: db,
	}, nil
}

func (e *ExplDB) Close() error {
	return e.db.Close()
}

func waitUntilAvailable(db *sql.DB) {
	for db.Ping() != nil {
		logrus.Info("Waiting for database...")
		time.Sleep(time.Second)
	}
}

//go:embed migrations/*.sql
var fs embed.FS

func applyMigrations(databaseUrl string) (err error) {
	handleDeferredCloseError := func(c io.Closer) {
		closeErr := c.Close()
		if closeErr != nil {
			err = multierror.Append(err, closeErr)
		}
	}

	srcDrv, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}
	defer handleDeferredCloseError(srcDrv)

	mig, err := migrate.NewWithSourceInstance("iofs", srcDrv, databaseUrl)
	if err != nil {
		return err
	}
	mig.Log = &migrateLoggerAdapter{}

	err = mig.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

type migrateLoggerAdapter struct {
}

func (r *migrateLoggerAdapter) Printf(format string, v ...interface{}) {
	logrus.Infof(format, v...)
}

func (r *migrateLoggerAdapter) Verbose() bool {
	return true
}
