package webhook

import (
	"fmt"
	"klio/expl/expldb"
	"klio/expl/types"
	"net/http"
	"regexp"
	"time"
	"unicode/utf16"
)

func NewAddHandler(edb expldb.Adder, settings AddHandlerSettings) Handler {
	return &addHandler{
		edb:      edb,
		settings: settings,
	}
}

type AddHandlerSettings interface {
	types.EntrySettings
	MaxUTF16LengthForKey() int
	MaxUTF16LengthForValue() int
}

type addHandler struct {
	edb      expldb.Adder
	settings AddHandlerSettings
}

func (a *addHandler) Handle(in *Request, r *http.Request, now time.Time) (*Response, error) {
	sep := regexp.MustCompile("^\\pZ*\\PZ+\\pZ+(?P<Key>\\PZ+)\\pZ+(?P<Value>\\PZ.*?)\\pZ*$")
	match := sep.FindStringSubmatch(in.Text)
	if match == nil {
		return NewResponse(fmt.Sprintf("Syntax: %s <Begriff> <Erklärung>", in.TriggerWord)), nil
	}
	key := match[sep.SubexpIndex("Key")]
	value := match[sep.SubexpIndex("Value")]

	if len(utf16.Encode([]rune(key))) > a.settings.MaxUTF16LengthForKey() {
		return NewResponse("Tut mir leid, der Begriff ist leider zu lang."), nil
	}

	if len(utf16.Encode([]rune(value))) > a.settings.MaxUTF16LengthForValue() {
		return NewResponse("Tut mir leid, die Erklärung ist leider zu lang."), nil
	}

	entry, err := a.edb.Add(r.Context(), key, value, in.UserName, now)
	if err != nil {
		return nil, err
	}

	return NewResponse(fmt.Sprintf("Ich habe den folgenden neuen Eintrag hinzugefügt:\n```\n%s\n```", entry.String(a.settings))), nil
}
