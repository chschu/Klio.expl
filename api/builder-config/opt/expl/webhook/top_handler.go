package webhook

import (
	"fmt"
	"klio/expl/expldb"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func NewTopHandler(edb expldb.Topper, settings TopHandlerSettings) Handler {
	return &topHandler{
		edb:      edb,
		settings: settings,
	}
}

type TopHandlerSettings interface {
	MaxTopCount() int
}

type topHandler struct {
	edb      expldb.Topper
	settings TopHandlerSettings
}

func (e *topHandler) Handle(in *Request, r *http.Request, _ time.Time) (*Response, error) {
	sep := regexp.MustCompile("^\\pZ*\\PZ+\\pZ*$")
	if !sep.MatchString(in.Text) {
		return NewResponse(fmt.Sprintf("Syntax: %s", in.TriggerWord)), nil
	}

	entries, err := e.edb.Top(r.Context(), e.settings.MaxTopCount())
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
