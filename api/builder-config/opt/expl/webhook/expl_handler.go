package webhook

import (
	"fmt"
	"klio/expl/expldb"
	"klio/expl/security"
	"klio/expl/types"
	"net/http"
	"regexp"
	"time"
)

func NewExplHandler(edb expldb.Explainer, indexSpecParser types.IndexSpecParser, entryListStringer EntryListStringer, maxImmediateResults int, webUrlPathPrefix string, webUrlValidity time.Duration) Handler {
	return &explHandler{
		edb:                 edb,
		entryListStringer:   entryListStringer,
		indexSpecParser:     indexSpecParser,
		maxImmediateResults: maxImmediateResults,
		webUrlPathPrefix:    webUrlPathPrefix,
		webUrlValidity:      webUrlValidity,
	}
}

type explHandler struct {
	edb               expldb.Explainer
	indexSpecParser   types.IndexSpecParser
	jwtGenerator      security.JwtGenerator
	entryListStringer EntryListStringer

	maxImmediateResults int
	webUrlPathPrefix    string
	webUrlValidity      time.Duration
}

func (e *explHandler) Handle(in *Request, r *http.Request, now time.Time) (*Response, error) {
	syntaxResponse := NewResponse(fmt.Sprintf("Syntax: %s <Begriff> ( <Index> | <VonIndex>:<BisIndex> )*", in.TriggerWord))

	sep := regexp.MustCompile("^\\pZ*\\PZ+\\pZ+(?P<Key>\\PZ+)(?:\\pZ+(?P<IndexSpec>.*?))?\\pZ*$")
	match := sep.FindStringSubmatch(in.Text)
	if match == nil {
		return syntaxResponse, nil
	}
	key := match[sep.SubexpIndex("Key")]
	indexSpecStr := match[sep.SubexpIndex("IndexSpec")]

	var indexSpec types.IndexSpec
	if len(indexSpecStr) == 0 {
		indexSpec = types.IndexSpecAll()
	} else {
		var err error
		indexSpec, err = e.indexSpecParser.ParseIndexSpec(indexSpecStr)
		if err != nil {
			return syntaxResponse, nil
		}
	}

	entries, total, err := e.edb.ExplainWithLimit(r.Context(), key, indexSpec, e.maxImmediateResults)
	if err != nil {
		return nil, err
	}

	text := e.entryListStringer.String(entries, total, key, e.webUrlPathPrefix, now.Add(e.webUrlValidity), r)

	return NewResponse(text), nil
}
