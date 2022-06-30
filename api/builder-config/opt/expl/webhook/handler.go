package webhook

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler interface {
	Handle(in *Request, r *http.Request) (*Response, error)
}

func NewHandlerAdapter(handler Handler, token string) http.Handler {
	return &handlerAdapter{
		handler: handler,
		token:   token,
	}
}

type handlerAdapter struct {
	handler Handler
	token   string
}

func (a *handlerAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	in := Request{}
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("error decoding request: %v", err)
		return
	}

	if in.Token != a.token {
		w.WriteHeader(http.StatusUnauthorized)
		logrus.Errorf("invalid token: %s", in.Token)
		return
	}

	out, err := a.handler.Handle(&in, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("error handling decoded request: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("error encoding and writing response: %v", err)
		return
	}
}
