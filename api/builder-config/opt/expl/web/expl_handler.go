package web

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"klio/expl/expldb"
	"klio/expl/types"
	"net/http"
)

func NewExplHandler(edb *expldb.ExplDB) http.Handler {
	return &explHandler{
		edb: edb,
	}
}

type explHandler struct {
	edb *expldb.ExplDB
}

func (e *explHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	entries, err := e.edb.Expl(key, types.IndexSpecAll())
	if err != nil {
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
			text = "Ich habe folgende Einträge gefunden:\n"
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
