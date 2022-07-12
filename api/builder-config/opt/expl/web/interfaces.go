package web

import (
	"github.com/chschu/Klio.expl/types"
	"net/http"
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source interfaces.go -destination generated/mocks/interfaces.go -package mocks

type EntryStringer interface {
	String(entry *types.Entry) string
}

type EntryListStringer interface {
	String(entries []types.Entry) string
}

type JWTValidator interface {
	ValidateJWT(jwtStr string) (subject string, err error)
}

type JWTExtractor interface {
	ExtractJWT(r *http.Request) string
}
