package webhook

import (
	"fmt"
	"klio/expl/expldb"
	"klio/expl/security"
	"net/http"
	"regexp"
	"time"
)

func NewFindHandler(edb expldb.Finder, entryListStringer EntryListStringer, maxImmediateResults int, webUrlPathPrefix string, webUrlValidity time.Duration) Handler {
	return &findHandler{
		edb:                 edb,
		entryListStringer:   entryListStringer,
		maxImmediateResults: maxImmediateResults,
		webUrlPathPrefix:    webUrlPathPrefix,
		webUrlValidity:      webUrlValidity,
	}
}

type findHandler struct {
	edb               expldb.Finder
	jwtGenerator      security.JwtGenerator
	entryListStringer EntryListStringer

	maxImmediateResults int
	webUrlPathPrefix    string
	webUrlValidity      time.Duration
}

func (f *findHandler) Handle(in *Request, r *http.Request, now time.Time) (*Response, error) {
	syntaxResponse := NewResponse(fmt.Sprintf("Syntax: %s <POSIX-Regex>", in.TriggerWord))

	sep := regexp.MustCompile("^\\pZ*\\PZ+\\pZ+(?P<Regex>\\PZ+)\\pZ*$")
	match := sep.FindStringSubmatch(in.Text)
	if match == nil {
		return syntaxResponse, nil
	}
	rex := match[sep.SubexpIndex("Regex")]

	entries, total, err := f.edb.FindWithLimit(r.Context(), rex, f.maxImmediateResults)
	if err == expldb.ErrFindRegexInvalid {
		return syntaxResponse, nil
	}
	if err != nil {
		return nil, err
	}

	text := f.entryListStringer.String(entries, total, rex, f.webUrlPathPrefix, now.Add(f.webUrlValidity), r)

	return NewResponse(text), nil
}
