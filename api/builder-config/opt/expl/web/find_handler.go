package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"klio/expl/expldb"
	"klio/expl/security"
	"net/http"
)

func NewFindHandler(edb *expldb.ExplDB, jwtValidate security.JwtValidate) http.Handler {
	return &findHandler{
		edb:         edb,
		jwtValidate: jwtValidate,
	}
}

type findHandler struct {
	edb         *expldb.ExplDB
	jwtValidate security.JwtValidate
}

func (f *findHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jwtStr := mux.Vars(r)["jwt"]

	rex, err := f.jwtValidate(jwtStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Infof("failed to validate JWT: %v", err)
		return
	}

	entries, err := f.edb.Find(rex)
	if err != nil && err != expldb.ErrFindRegexInvalid {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("error accessing entries: %v", err)
		return
	}

	count := len(entries)

	var text string
	if count == 0 {
		text = "Ich habe leider keinen Eintrag gefunden.\n"
	} else {
		if count == 1 {
			text = "Ich habe den folgenden Eintrag gefunden:\n"
		} else {
			text = fmt.Sprintf("Ich habe die folgenden %d Einträge gefunden:\n", count)
		}
		for _, entry := range entries {
			text += entry.String() + "\n"
		}
	}

	h := w.Header()
	h.Set("Content-Type", "text/plain; charset=UTF-8")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Pragma", "no-cache")
	h.Set("Expires", "0")
	h.Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")

	_, err = w.Write([]byte(text))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("error writing response: %v", err)
		return
	}
}
