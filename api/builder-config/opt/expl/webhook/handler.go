package webhook

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler interface {
	Handle(in *Request, r *http.Request) (*Response, error)
}

type HandlerFunc func(in *Request, r *http.Request) (*Response, error)

func (f HandlerFunc) Handle(in *Request, r *http.Request) (*Response, error) {
	return f(in, r)
}

func RequiredTokenAdapter(token string) func(handler Handler) Handler {
	return func(handler Handler) Handler {
		return HandlerFunc(func(in *Request, r *http.Request) (*Response, error) {
			if in.Token != token {
				return nil, fmt.Errorf("invalid token: %s", in.Token)
			}
			return handler.Handle(in, r)
		})
	}
}

func ToHttpHandler(handler Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in := Request{}
		err := json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.Errorf("error decoding request: %v", err)
			return
		}

		out, err := handler.Handle(&in, r)
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
	})
}
