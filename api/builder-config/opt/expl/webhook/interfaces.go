package webhook

import (
	"github.com/chschu/Klio.expl/types"
	"net/http"
	"time"
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source interfaces.go -destination generated/mocks/interfaces.go -package mocks

type EntryStringer interface {
	String(entry *types.Entry) string
}

type EntryListStringer interface {
	String(entries []types.Entry, total int, subject string, webUrlPathPrefix string, webUrlExpiresAt time.Time, req *http.Request) string
}

type JWTGenerator interface {
	GenerateJWT(subject string, expiresAt time.Time) (jwtStr string, err error)
}
