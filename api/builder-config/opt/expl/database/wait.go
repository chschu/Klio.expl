package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"
)

func WaitUntilAvailable(db *sqlx.DB) {
	for db.Ping() != nil {
		logrus.Info("Waiting for database...")
		time.Sleep(time.Second)
	}
}
