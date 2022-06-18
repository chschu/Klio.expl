package main

import (
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"klio/expl/expldb"
	"klio/expl/webhook"
	"net/http"
	"os"
	"time"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})

	edb, err := expldb.Init(mustLookupEnv("CONNECT_STRING"))
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

	http.Handle("/expl/add", webhook.NewAddHandler(edb, mustLookupEnv("WEBHOOK_TOKEN_ADD")))

	logrus.Info("Listening for HTTP connections...")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Shutting down")
}

func mustLookupEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		logrus.Fatalf("environment variable not set: %v", key)
	}
	return value
}
