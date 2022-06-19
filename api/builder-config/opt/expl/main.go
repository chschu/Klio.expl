package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"klio/expl/expldb"
	"klio/expl/web"
	"klio/expl/webhook"
	"net/http"
	"os"
	"strconv"
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

	var wrap func(http.Handler) http.Handler
	if mustLookupUseProxyHeaders() {
		wrap = handlers.ProxyHeaders
	} else {
		wrap = func(h http.Handler) http.Handler { return h }
	}

	r := mux.NewRouter()
	r.Handle("/api/add", wrap(webhook.NewAddHandler(edb, mustLookupEnv("WEBHOOK_TOKEN_ADD"))))
	r.Handle("/api/expl", wrap(webhook.NewExplHandler(edb, mustLookupEnv("WEBHOOK_TOKEN_EXPL"), "/expl/")))
	r.Handle("/api/del", wrap(webhook.NewDelHandler(edb, mustLookupEnv("WEBHOOK_TOKEN_DEL"))))
	r.Handle("/expl/{key:.*}", wrap(web.NewExplHandler(edb)))

	logrus.Info("Listening for HTTP connections...")
	err = http.ListenAndServe(":8000", r)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Shutting down")
}

func mustLookupUseProxyHeaders() bool {
	envStr := mustLookupEnv("USE_PROXY_HEADERS")
	result, err := strconv.ParseBool(envStr)
	if err != nil {
		logrus.Fatalf("cannot convert value: %v", err)
	}
	return result
}

func mustLookupEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		logrus.Fatalf("environment variable not set: %v", key)
	}
	return value
}
