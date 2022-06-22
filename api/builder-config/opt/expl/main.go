package main

import (
	"context"
	"crypto/rand"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"klio/expl/expldb"
	"klio/expl/importer"
	"klio/expl/security"
	"klio/expl/settings"
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

	jwtKey := make([]byte, 256/8)
	_, err := rand.Read(jwtKey)
	if err != nil {
		logrus.Fatal(err)
	}
	jwtGenerator, jwtValidator := security.NewJwtHandlers(jwt.SigningMethodHS256, jwtKey, jwtKey)

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

	err = importer.Import()
	if err != nil {
		logrus.Fatal(err)
	}

	wrap := timeoutHandlerTransform(settings.HandlerTimeout)
	if mustLookupUseProxyHeaders() {
		wrap = wrap.compose(handlers.ProxyHeaders)
	}

	r := mux.NewRouter()
	r.Handle("/api/add", wrap(webhook.NewAddHandler(edb, mustLookupEnv("WEBHOOK_TOKEN_ADD"))))
	r.Handle("/api/expl", wrap(webhook.NewExplHandler(edb, mustLookupEnv("WEBHOOK_TOKEN_EXPL"), "/expl/", jwtGenerator)))
	r.Handle("/api/del", wrap(webhook.NewDelHandler(edb, mustLookupEnv("WEBHOOK_TOKEN_DEL"))))
	r.Handle("/api/find", wrap(webhook.NewFindHandler(edb, mustLookupEnv("WEBHOOK_TOKEN_FIND"), "/find/", jwtGenerator)))
	r.Handle("/api/top", wrap(webhook.NewTopHandler(edb, mustLookupEnv("WEBHOOK_TOKEN_TOP"))))
	r.Handle("/expl/{jwt:.*}", wrap(web.NewExplHandler(edb, jwtValidator)))
	r.Handle("/find/{jwt:.*}", wrap(web.NewFindHandler(edb, jwtValidator)))

	logrus.Info("Listening for HTTP connections...")
	err = http.ListenAndServe(":8000", r)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Shutting down")
}

type handlerTransform func(http.Handler) http.Handler

func (outer handlerTransform) compose(inner handlerTransform) handlerTransform {
	return func(h http.Handler) http.Handler {
		return outer(inner(h))
	}
}

func timeoutHandlerTransform(timeout time.Duration) handlerTransform {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
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
