package main

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"klio/expl/expldb"
	"klio/expl/web"
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

	r := mux.NewRouter()
	r.Handle("/api/add", webhook.NewAddHandler(edb, mustLookupEnv("WEBHOOK_TOKEN_ADD")))
	r.Handle("/api/expl", webhook.NewExplHandler(edb, mustLookupEnv("WEBHOOK_TOKEN_EXPL"), "/expl/"))
	r.Handle("/api/del", webhook.NewDelHandler(edb, mustLookupEnv("WEBHOOK_TOKEN_DEL")))
	r.Handle("/expl/{key:.*}", web.NewExplHandler(edb))

	logrus.Info("Listening for HTTP connections...")
	err = http.ListenAndServe(":8000", r)
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
