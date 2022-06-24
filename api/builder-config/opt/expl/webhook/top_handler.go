package webhook

import (
	"fmt"
	"klio/expl/expldb"
	"klio/expl/settings"
	"net/http"
	"regexp"
	"strings"
)

func NewTopHandler(edb *expldb.ExplDB, token string) http.Handler {
	return NewHandlerAdapter(&topHandler{
		edb:   edb,
		token: token,
	})
}

type topHandler struct {
	edb   *expldb.ExplDB
	token string
}

func (t *topHandler) Token() string {
	return t.token
}

func (e *topHandler) Handle(in *Request, r *http.Request) (*Response, error) {
	sep := regexp.MustCompile("^\\pZ*\\PZ+\\pZ*$")
	if !sep.MatchString(in.Text) {
		return NewResponse(fmt.Sprintf("Syntax: %s", in.TriggerWord)), nil
	}

	entries, err := e.edb.Top(settings.TopExplCount)
	if err != nil {
		return nil, err
	}

	count := len(entries)

	var text string
	if count == 0 {
		text = "Ich habe leider keinen Eintrag gefunden."
	} else {
		parts := make([]string, 0, count)
		for _, entry := range entries {
			parts = append(parts, fmt.Sprintf("%s(%d)", entry.Key, entry.HeadIndex))
		}
		if count == 1 {
			text = "Ich habe den folgenden wichtigsten Eintrag gefunden:\n"
		} else {
			text = fmt.Sprintf("Ich habe die folgenden %d wichtigsten Einträge gefunden:\n", count)
		}

		text += "```\n" + strings.Join(parts, ", ") + "\n```"
	}

	return NewResponse(text), nil
}
