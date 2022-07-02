package webhook

import (
	"fmt"
	"klio/expl/expldb"
	"klio/expl/types"
	"klio/expl/util"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func NewDelHandler(edb expldb.Deleter, entryStringer util.EntryStringer) Handler {
	return &delHandler{
		edb:           edb,
		entryStringer: entryStringer,
	}
}

type delHandler struct {
	edb           expldb.Deleter
	entryStringer util.EntryStringer
}

func (d *delHandler) Handle(in *Request, r *http.Request, _ time.Time) (*Response, error) {
	syntaxResponse := NewResponse(fmt.Sprintf("Syntax: %s <Begriff> <Index>", in.TriggerWord))

	sep := regexp.MustCompile("^\\pZ*\\PZ+\\pZ+(?P<Key>\\PZ+)\\pZ+(?P<Index>.*?)\\pZ*$")
	match := sep.FindStringSubmatch(in.Text)
	if match == nil {
		return syntaxResponse, nil
	}
	key := match[sep.SubexpIndex("Key")]
	indexStr := match[sep.SubexpIndex("Index")]

	index, err := parseIndex(indexStr)
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
