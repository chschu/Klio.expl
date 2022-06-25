package webhook

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"klio/expl/expldb"
	"klio/expl/security"
	"klio/expl/settings"
	"klio/expl/types"
	"net/http"
	"net/url"
	"regexp"
)

func NewFindHandler(edb *expldb.ExplDB, token string, webFindPathPrefix string, jwtGenerate security.JwtGenerate) http.Handler {
	return NewHandlerAdapter(&findHandler{
		edb:               edb,
		token:             token,
		webFindPathPrefix: webFindPathPrefix,
		jwtGenerate:       jwtGenerate,
	})
}

type findHandler struct {
	edb               *expldb.ExplDB
	token             string
	webFindPathPrefix string
	jwtGenerate       security.JwtGenerate
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
			var entriesText string
			findUrl, err := f.getWebFindUrl(r, rex)
			if err != nil {
				logrus.Warnf("unable to resolve URL for web find: %v", err)
				entriesText = fmt.Sprintf("%d Einträge", count)
			} else {
				entriesText = fmt.Sprintf("[%d Einträge](%s)", count, findUrl)
			}
			text = fmt.Sprintf("Ich habe %s gefunden, das sind die letzten %d:\n",
				entriesText, settings.MaxFindCount)
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

func (f *findHandler) getWebFindUrl(r *http.Request, rex string) (*url.URL, error) {
	scheme := r.URL.Scheme
	if scheme == "" {
		if r.TLS == nil {
			scheme = "http"
		} else {
			scheme = "https"
		}
	}
	jwtStr, err := f.jwtGenerate(rex, settings.FindTokenValidity)
	if err != nil {
		return nil, err
	}
	return url.Parse(scheme + "://" + r.Host + f.webFindPathPrefix + jwtStr)
}
