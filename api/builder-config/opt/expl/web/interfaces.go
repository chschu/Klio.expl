package web

import "klio/expl/types"

type EntryStringer interface {
	String(entry *types.Entry) string
}

type JwtValidator interface {
	Validate(jwtStr string) (subject string, err error)
}
