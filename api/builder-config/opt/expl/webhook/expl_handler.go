package webhook

import (
	"context"
	"fmt"
	"klio/expl/types"
	"net/http"
	"regexp"
	"time"
)

type IndexSpecParser interface {
	ParseIndexSpec(s string) (types.IndexSpec, error)
}

type LimitedExplainer interface {
	ExplainWithLimit(ctx context.Context, key string, indexSpec types.IndexSpec, limit int) (entries []types.Entry, total int, err error)
}

func NewExplHandler(edb LimitedExplainer, indexSpecParser IndexSpecParser, entryListStringer EntryListStringer, maxImmediateResults int, webUrlPathPrefix string, webUrlValidity time.Duration) *explHandler {
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
	edb               LimitedExplainer
	indexSpecParser   IndexSpecParser
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
