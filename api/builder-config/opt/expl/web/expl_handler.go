package web

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"klio/expl/types"
	"net/http"
)

type Explainer interface {
	Explain(ctx context.Context, key string, indexSpec types.IndexSpec) (entries []types.Entry, err error)
}

func NewExplHandler(edb Explainer, jwtValidate JwtValidator, entryListStringer EntryListStringer) *explHandler {
	return &explHandler{
		edb:               edb,
		jwtValidator:      jwtValidate,
		entryListStringer: entryListStringer,
	}
}

type explHandler struct {
	edb               Explainer
	jwtValidator      JwtValidator
	entryListStringer EntryListStringer
}

func (e *explHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jwtStr := mux.Vars(r)["jwt"]

	key, err := e.jwtValidator.Validate(jwtStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Infof("failed to validate JWT: %v", err)
		return
	}

	entries, err := e.edb.Explain(r.Context(), key, types.IndexSpecAll())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("error accessing entries: %v", err)
		return
	}

	h := w.Header()
	h.Set("Content-Type", "text/plain; charset=UTF-8")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Pragma", "no-cache")
	h.Set("Expires", "0")
	h.Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")

	_, err = io.WriteString(w, e.entryListStringer.String(entries))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("error writing response: %v", err)
		return
	}
}
