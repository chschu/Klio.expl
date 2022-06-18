package webhook

import (
	"fmt"
	"golang.org/x/text/unicode/norm"
	"klio/expl/expldb"
	"klio/expl/settings"
	"net/http"
	"regexp"
	"time"
	"unicode/utf8"
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

func (a *addHandler) Handle(in *Request) (*Response, error) {
	sep := regexp.MustCompile("^\\pZ*\\PZ+\\pZ+(?P<Key>\\PZ+)\\pZ+(?P<Value>\\PZ.*?)\\pZ*$")
	match := sep.FindStringSubmatch(in.Text)
	if match == nil {
		return NewResponse(fmt.Sprintf("Syntax: %s <Begriff> <Erklärung>", in.TriggerWord)), nil
	}
	key := match[sep.SubexpIndex("Key")]
	value := match[sep.SubexpIndex("Value")]

	if utf8.RuneCountInString(norm.NFC.String(key)) > settings.MaxRuneCountForNormalizedKey {
		return NewResponse("Tut mir leid, der Begriff ist leider zu lang."), nil
	}

	if utf8.RuneCountInString(norm.NFC.String(value)) > settings.MaxRuneCountForNormalizedValue {
		return NewResponse("Tut mir leid, die Erklärung ist leider zu lang."), nil
	}

	entry, err := a.edb.Add(key, value, in.UserName, time.Now())
	if err != nil {
		return nil, err
	}

	return NewResponse(fmt.Sprintf(
		"Ich habe den folgenden neuen Eintrag hinzugefügt:\n```\n%s[%s/%s]\n```",
		entry.Key,
		entry.HeadIndex,
		entry.PermanentIndex)), nil
}
