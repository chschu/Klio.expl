package webhook

import (
	"fmt"
	"klio/expl/expldb"
	"klio/expl/types"
	"net/http"
	"regexp"
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

func (d *delHandler) Handle(in *Request) (*Response, error) {
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

	entries, err := d.edb.Del(key, types.IndexSpecSingle(index))
	if err != nil {
		return nil, err
	}

	var entriesText = ""
	for _, entry := range entries {
		entriesText += entry.String() + "\n"
	}

	var text string
	switch len(entries) {
	case 0:
		text = "Ich habe leider keinen Eintrag zum Löschen gefunden."
	case 1:
		text = fmt.Sprintf("Ich habe den folgenden Eintrag gelöscht:\n```\n%s```", entriesText)
	default:
		text = fmt.Sprintf("Ich habe folgende Einträge gelöscht:\n```\n%s```", entriesText)
	}

	return NewResponse(text), nil
}
