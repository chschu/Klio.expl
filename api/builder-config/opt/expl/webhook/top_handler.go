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

type Topper interface {
	Top(ctx context.Context, count int) (entries []types.Entry, err error)
}

func NewTopHandler(edb Topper, maxResults int) *topHandler {
	return &topHandler{
		edb:        edb,
		maxResults: maxResults,
	}
}

type topHandler struct {
	edb Topper

	maxResults int
}

var topRegexp = regexp.MustCompile("^\\pZ*\\PZ+\\pZ*$")

func (e *topHandler) Handle(in *Request, r *http.Request, _ time.Time) (*Response, error) {
	if !topRegexp.MatchString(in.Text) {
		return NewResponse(fmt.Sprintf("Syntax: %s", in.TriggerWord)), nil
	}

	entries, err := e.edb.Top(r.Context(), e.maxResults)
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
