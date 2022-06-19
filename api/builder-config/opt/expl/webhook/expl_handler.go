package webhook

import (
	"fmt"
	"klio/expl/expldb"
	"klio/expl/settings"
	"klio/expl/types"
	"net/http"
	"regexp"
)

func NewExplHandler(edb *expldb.ExplDB, token string) http.Handler {
	return NewHandlerAdapter(&explHandler{
		edb:   edb,
		token: token,
	})
}

type explHandler struct {
	edb   *expldb.ExplDB
	token string
}

func (e *explHandler) Token() string {
	return e.token
}

func (e *explHandler) Handle(in *Request) (*Response, error) {
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
		indexSpec = []types.IndexRange{{types.HeadIndex(1), types.TailIndex(1)}}
	} else {
		var err error
		indexSpec, err = parseIndexSpec(indexSpecStr)
		if err != nil {
			return syntaxResponse, nil
		}
	}

	entries, err := e.edb.Expl(key, indexSpec)
	if err != nil {
		return nil, err
	}

	count := len(entries)

	var text string
	if count == 0 {
		text = "Ich habe leider keinen Eintrag gefunden."
	} else {
		var caption string
		var limitedEntries []types.Entry
		if count > settings.MaxExplCount {
			// TODO add URL
			caption = fmt.Sprintf("Ich habe %d Einträge gefunden, das sind die letzten %d:",
				count, settings.MaxExplCount)
			limitedEntries = entries[count-settings.MaxExplCount:]
		} else {
			if count == 1 {
				caption = "Ich habe den folgenden Eintrag gefunden:"
			} else {
				caption = "Ich habe folgende Einträge gefunden:"
			}
			limitedEntries = entries
		}

		var entriesText = ""
		for _, entry := range limitedEntries {
			entriesText += entry.String() + "\n"
		}

		text = fmt.Sprintf("%s\n```\n%s```", caption, entriesText)
	}

	return NewResponse(text), nil
}
