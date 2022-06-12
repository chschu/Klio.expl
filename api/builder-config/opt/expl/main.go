package main

import (
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"klio/expl/expldb"
	"net/http"
	"os"
	"time"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})

	edb, err := expldb.Init(os.Getenv("CONNECT_STRING"))
	if err != nil {
		logrus.Fatal(err)
	}
	defer func(e *expldb.ExplDB) {
		err := e.Close()
		if err != nil {
			logrus.Fatal(err)
		}
	}(edb)
	logrus.Info("Database successfully initialized")

	logrus.Info("Listening for HTTP connections...")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Shutting down")
}
