package web

import "klio/expl/types"

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source interfaces.go -destination generated/mocks/interfaces.go -package mocks

type EntryStringer interface {
	String(entry *types.Entry) string
}

type EntryListStringer interface {
	String(entries []types.Entry) string
}

type JwtValidator interface {
	Validate(jwtStr string) (subject string, err error)
}
