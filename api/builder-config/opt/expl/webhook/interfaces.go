package webhook

import (
	"klio/expl/types"
	"net/http"
	"time"
)

type EntryStringer interface {
	String(entry *types.Entry) string
}

type EntryListStringer interface {
	String(entries []types.Entry, total int, subject string, webUrlPathPrefix string, webUrlExpiresAt time.Time, req *http.Request) string
}

type JwtGenerator interface {
	Generate(subject string, expiresAt time.Time) (jwtStr string, err error)
}
