package database

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"syscall"
	"time"
)

func WaitUntilAvailable(ctx context.Context, db *sqlx.DB) error {
	for {
		err := db.PingContext(ctx)
		if err == nil || !errors.Is(err, syscall.ECONNREFUSED) {
			return err
		}
		logrus.Info("Waiting for database...")
		time.Sleep(time.Second)
	}
}
