package webhook

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"klio/expl/expldb"
	"klio/expl/settings"
	"klio/expl/types"
	"net/http"
	"net/url"
	"regexp"
)

func NewFindHandler(edb *expldb.ExplDB, token string, webFindPathPrefix string) http.Handler {
	return NewHandlerAdapter(&findHandler{
		edb:               edb,
		token:             token,
		webFindPathPrefix: webFindPathPrefix,
	})
}

type findHandler struct {
	edb               *expldb.ExplDB
	token             string
	webFindPathPrefix string
}

func (f *findHandler) Token() string {
	return f.token
}

func (f *findHandler) Handle(in *Request, r *http.Request) (*Response, error) {
	syntaxResponse := NewResponse(fmt.Sprintf("Syntax: %s <POSIX-Regex>", in.TriggerWord))

	sep := regexp.MustCompile("^\\pZ*\\PZ+\\pZ+(?P<Regex>\\PZ+)\\pZ*$")
	match := sep.FindStringSubmatch(in.Text)
	if match == nil {
		return syntaxResponse, nil
	}
	rex := match[sep.SubexpIndex("Regex")]

	entries, err := f.edb.Find(rex)
	if err == expldb.ErrFindRegexInvalid {
		return syntaxResponse, nil
	}
	if err != nil {
		return nil, err
	}

	count := len(entries)

	var text string
	if count == 0 {
		text = "Ich habe leider keinen Eintrag gefunden."
	} else {
		var limitedEntries []types.Entry
		if count > settings.MaxFindCount {
			urlText := ""
			findUrl, err := f.getWebFindUrl(r, rex)
			if err != nil {
				logrus.Warnf("unable to resolve URL for web expl: %v", err)
			} else {
				urlText = fmt.Sprintf(" (%s)", findUrl)
			}
			text = fmt.Sprintf("Ich habe %d Einträge gefunden%s, das sind die letzten %d:\n",
				count, urlText, settings.MaxFindCount)
			limitedEntries = entries[count-settings.MaxFindCount:]
		} else {
			if count == 1 {
				text = "Ich habe den folgenden Eintrag gefunden:\n"
			} else {
				text = fmt.Sprintf("Ich habe die folgenden %d Einträge gefunden:\n", count)
			}
			limitedEntries = entries
		}
		text += "```\n"
		for _, entry := range limitedEntries {
			text += entry.String() + "\n"
		}
		text += "```"
	}

	return NewResponse(text), nil
}

func (f *findHandler) getWebFindUrl(r *http.Request, regexStr string) (*url.URL, error) {
	scheme := r.URL.Scheme
	if scheme == "" {
		if r.TLS == nil {
			scheme = "http"
		} else {
			scheme = "https"
		}
	}
	return url.Parse(scheme + "://" + r.Host + f.webFindPathPrefix + url.PathEscape(regexStr))
}
