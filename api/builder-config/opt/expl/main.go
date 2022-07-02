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

	edb, err := expldb.NewExplDB(mustLookupEnv("CONNECT_STRING"))
	if err != nil {
		logrus.Fatal(err)
	}
	defer func(e expldb.ExplDB) {
		err := e.Close()
		if err != nil {
			logrus.Fatal(err)
		}
	}(edb)
	logrus.Info("Database successfully initialized")

	useProxyHeaders := mustParseBool(mustLookupEnv("USE_PROXY_HEADERS"))
	webChain := compose(timeoutAdapter(settings.Instance.HandlerTimeout()), proxyHeaderAdapter(useProxyHeaders))
	webhookChain := compose(webChain, webhook.ToHttpHandler)

	r := mux.NewRouter()
	r.Handle("/api/add", compose(webhookChain, requiredTokenEnvAdapter("WEBHOOK_TOKEN_ADD"))(webhook.NewAddHandler(edb, settings.Instance)))
	r.Handle("/api/expl", compose(webhookChain, requiredTokenEnvAdapter("WEBHOOK_TOKEN_EXPL"))(webhook.NewExplHandler(edb, "/expl/", jwtGenerator, settings.Instance)))
	r.Handle("/api/del", compose(webhookChain, requiredTokenEnvAdapter("WEBHOOK_TOKEN_DEL"))(webhook.NewDelHandler(edb, settings.Instance)))
	r.Handle("/api/find", compose(webhookChain, requiredTokenEnvAdapter("WEBHOOK_TOKEN_FIND"))(webhook.NewFindHandler(edb, "/find/", jwtGenerator, settings.Instance)))
	r.Handle("/api/top", compose(webhookChain, requiredTokenEnvAdapter("WEBHOOK_TOKEN_TOP"))(webhook.NewTopHandler(edb, settings.Instance)))
	r.Handle("/expl/{jwt:.*}", webChain(web.NewExplHandler(edb, jwtValidator, settings.Instance)))
	r.Handle("/find/{jwt:.*}", webChain(web.NewFindHandler(edb, jwtValidator, settings.Instance)))

	logrus.Info("Listening for HTTP connections...")
	err = http.ListenAndServe(":8000", r)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Shutting down")
}

func proxyHeaderAdapter(useProxyHeaders bool) func(http.Handler) http.Handler {
	if useProxyHeaders {
		return handlers.ProxyHeaders
	}
	return func(handler http.Handler) http.Handler {
		return handler
	}
}

func requiredTokenEnvAdapter(envKey string) func(webhook.Handler) webhook.Handler {
	return webhook.RequiredTokenAdapter(mustLookupEnv(envKey))
}

func timeoutAdapter(timeout time.Duration) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func mustLookupEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		logrus.Fatalf("environment variable not set: %v", key)
	}
	return value
}

func mustParseBool(s string) bool {
	result, err := strconv.ParseBool(s)
	if err != nil {
		logrus.Fatalf("cannot convert to bool: %s", s)
	}
	return result
}

func compose[A any, B any, C any](outer func(B) C, inner func(A) B) func(A) C {
	return func(a A) C {
		return outer(inner(a))
	}
}
