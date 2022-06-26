package webhook

import (
	"fmt"
	"klio/expl/expldb"
	"klio/expl/settings"
	"net/http"
	"regexp"
	"time"
	"unicode/utf16"
)

func NewAddHandler(edb *expldb.ExplDB, token string) http.Handler {
	return NewHandlerAdapter(&addHandler{
		edb:   edb,
		token: token,
	})
}

type addHandler struct {
	edb   *expldb.ExplDB
	token string
}

func (a *addHandler) Token() string {
	return a.token
}

func (a *addHandler) Handle(in *Request, r *http.Request) (*Response, error) {
	sep := regexp.MustCompile("^\\pZ*\\PZ+\\pZ+(?P<Key>\\PZ+)\\pZ+(?P<Value>\\PZ.*?)\\pZ*$")
	match := sep.FindStringSubmatch(in.Text)
	if match == nil {
		return NewResponse(fmt.Sprintf("Syntax: %s <Begriff> <Erklärung>", in.TriggerWord)), nil
	}
	key := match[sep.SubexpIndex("Key")]
	value := match[sep.SubexpIndex("Value")]

	if len(utf16.Encode([]rune(key))) > settings.MaxUTF16LengthForKey {
		return NewResponse("Tut mir leid, der Begriff ist leider zu lang."), nil
	}

	if len(utf16.Encode([]rune(value))) > settings.MaxUTF16LengthForValue {
		return NewResponse("Tut mir leid, die Erklärung ist leider zu lang."), nil
	}

	entry, err := a.edb.Add(r.Context(), key, value, in.UserName, time.Now())
	if err != nil {
		return nil, err
	}

	return NewResponse(fmt.Sprintf("Ich habe den folgenden neuen Eintrag hinzugefügt:\n```\n%s\n```", entry)), nil
}
