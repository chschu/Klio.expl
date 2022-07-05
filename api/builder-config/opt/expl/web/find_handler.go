package web

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"klio/expl/expldb"
	"klio/expl/types"
	"net/http"
)

type Finder interface {
	Find(ctx context.Context, rex string) (entries []types.Entry, err error)
}

func NewFindHandler(edb Finder, jwtValidator JWTValidator, entryListStringer EntryListStringer) *findHandler {
	return &findHandler{
		edb:               edb,
		jwtValidator:      jwtValidator,
		entryListStringer: entryListStringer,
	}
}

type findHandler struct {
	edb               Finder
	jwtValidator      JWTValidator
	entryListStringer EntryListStringer
}

func (f *findHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jwtStr := mux.Vars(r)["jwt"]

	rex, err := f.jwtValidator.ValidateJWT(jwtStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Infof("failed to validate JWT: %v", err)
		return
	}

	entries, err := f.edb.Find(r.Context(), rex)
	if err != nil && err != expldb.ErrFindRegexInvalid {
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

	_, err = io.WriteString(w, f.entryListStringer.String(entries))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("error writing response: %v", err)
		return
	}
}
