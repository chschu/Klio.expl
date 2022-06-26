package webhook

import (
	"fmt"
	"klio/expl/expldb"
	"klio/expl/types"
	"net/http"
	"regexp"
	"strings"
)

func NewDelHandler(edb *expldb.ExplDB, token string) http.Handler {
	return NewHandlerAdapter(&delHandler{
		edb:   edb,
		token: token,
	})
}

type delHandler struct {
	edb   *expldb.ExplDB
	token string
}

func (d *delHandler) Token() string {
	return d.token
}

func (d *delHandler) Handle(in *Request, r *http.Request) (*Response, error) {
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

	entries, err := d.edb.Del(r.Context(), key, types.IndexSpecSingle(index))
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
			sb.WriteString(entry.String())
			sb.WriteRune('\n')
		}
		sb.WriteString("```")
	}

	return NewResponse(sb.String()), nil
}
