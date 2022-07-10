package webhook

import (
	"context"
	"fmt"
	"klio/expl/types"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type IndexParser interface {
	ParseIndex(s string) (types.Index, error)
}

type Deleter interface {
	Delete(ctx context.Context, key string, indexSpec types.IndexSpec) (entries []types.Entry, err error)
}

func NewDelHandler(edb Deleter, indexParser IndexParser, entryStringer EntryStringer) *delHandler {
	return &delHandler{
		edb:           edb,
		indexParser:   indexParser,
		entryStringer: entryStringer,
	}
}

type delHandler struct {
	edb           Deleter
	indexParser   IndexParser
	entryStringer EntryStringer
}

var delHandlerRegexp = regexp.MustCompile("^\\pZ*\\PZ+\\pZ+(?P<Key>\\PZ+)\\pZ+(?P<Index>.*?)\\pZ*$")
var delHandlerSubexpIndexKey = delHandlerRegexp.SubexpIndex("Key")
var delHandlerSubexpIndexIndex = delHandlerRegexp.SubexpIndex("Index")

func (d *delHandler) Handle(in *Request, r *http.Request, _ time.Time) (*Response, error) {
	syntaxResponse := NewResponse(fmt.Sprintf("Syntax: %s <Begriff> <Index>", in.TriggerWord))

	match := delHandlerRegexp.FindStringSubmatch(in.Text)
	if match == nil {
		return syntaxResponse, nil
	}
	key := match[delHandlerSubexpIndexKey]
	indexStr := match[delHandlerSubexpIndexIndex]

	index, err := d.indexParser.ParseIndex(indexStr)
	if err != nil {
		return syntaxResponse, nil
	}

	entries, err := d.edb.Delete(r.Context(), key, types.IndexSpecSingle(index))
	if err != nil {
		return nil, err
	}

	count := len(entries)

	var sb strings.Builder
	if count == 0 {
		sb.WriteString("Ich habe leider keinen Eintrag zum Löschen gefunden.")
	} else {
		if count == 1 {
			sb.WriteString("Ich habe den folgenden Eintrag gelöscht:\n")
		} else {
			sb.WriteString(fmt.Sprintf("Ich habe die folgenden %d Einträge gelöscht:\n", count))
		}
		sb.WriteString("```\n")
		for _, entry := range entries {
			sb.WriteString(d.entryStringer.String(&entry))
			sb.WriteRune('\n')
		}
		sb.WriteString("```")
	}

	return NewResponse(sb.String()), nil
}
