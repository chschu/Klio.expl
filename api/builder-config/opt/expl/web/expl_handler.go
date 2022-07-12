package web

import (
	"context"
	"github.com/chschu/Klio.expl/types"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Explainer interface {
	Explain(ctx context.Context, key string) (entries []types.Entry, err error)
}

func NewExplHandler(edb Explainer, jwtExtractor JWTExtractor, jwtValidator JWTValidator, entryListStringer EntryListStringer) *explHandler {
	return &explHandler{
		edb:               edb,
		jwtExtractor:      jwtExtractor,
		jwtValidator:      jwtValidator,
		entryListStringer: entryListStringer,
	}
}

type explHandler struct {
	edb               Explainer
	jwtExtractor      JWTExtractor
	jwtValidator      JWTValidator
	entryListStringer EntryListStringer
}

func (e *explHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jwtStr := e.jwtExtractor.ExtractJWT(r)

	key, err := e.jwtValidator.ValidateJWT(jwtStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Infof("failed to validate JWT: %v", err)
		return
	}

	entries, err := e.edb.Explain(r.Context(), key)
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
