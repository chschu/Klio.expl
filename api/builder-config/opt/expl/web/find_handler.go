package web

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"klio/expl/expldb"
	"klio/expl/security"
	"net/http"
)

func NewFindHandler(edb *expldb.ExplDB, jwtValidator security.JwtValidator) http.Handler {
	return &findHandler{
		edb:          edb,
		jwtValidator: jwtValidator,
	}
}

type findHandler struct {
	edb          *expldb.ExplDB
	jwtValidator security.JwtValidator
}

func (f *findHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jwtStr := mux.Vars(r)["jwt"]

	rex, err := f.jwtValidator.Validate(jwtStr)
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

	count := len(entries)

	var buf bytes.Buffer
	if count == 0 {
		buf.WriteString("Ich habe leider keinen Eintrag gefunden.\n")
	} else {
		if count == 1 {
			buf.WriteString("Ich habe den folgenden Eintrag gefunden:\n")
		} else {
			buf.WriteString(fmt.Sprintf("Ich habe die folgenden %d Einträge gefunden:\n", count))
		}
		for _, entry := range entries {
			buf.WriteString(entry.String())
			buf.WriteRune('\n')
		}
	}

	h := w.Header()
	h.Set("Content-Type", "text/plain; charset=UTF-8")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Pragma", "no-cache")
	h.Set("Expires", "0")
	h.Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")

	_, err = buf.WriteTo(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("error writing response: %v", err)
		return
	}
}
