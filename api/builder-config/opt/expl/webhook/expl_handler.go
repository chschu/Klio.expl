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

func NewExplHandler(edb *expldb.ExplDB, token string, webExplPathPrefix string, jwtGenerate security.JwtGenerate) http.Handler {
	return NewHandlerAdapter(&explHandler{
		edb:               edb,
		token:             token,
		webExplPathPrefix: webExplPathPrefix,
		jwtGenerate:       jwtGenerate,
	})
}

type explHandler struct {
	edb               *expldb.ExplDB
	token             string
	webExplPathPrefix string
	jwtGenerate       security.JwtGenerate
}

func (e *explHandler) Token() string {
	return e.token
}

func (e *explHandler) Handle(in *Request, r *http.Request) (*Response, error) {
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
		indexSpec = types.IndexSpecAll()
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
		var limitedEntries []types.Entry
		if count > settings.MaxExplCount {
			urlText := ""
			explUrl, err := e.getWebExplUrl(r, key)
			if err != nil {
				logrus.Warnf("unable to resolve URL for web expl: %v", err)
			} else {
				urlText = fmt.Sprintf(" (%s)", explUrl)
			}
			text = fmt.Sprintf("Ich habe %d Einträge gefunden%s, das sind die letzten %d:\n",
				count, urlText, settings.MaxExplCount)
			limitedEntries = entries[count-settings.MaxExplCount:]
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

func (e *explHandler) getWebExplUrl(r *http.Request, key string) (*url.URL, error) {
	scheme := r.URL.Scheme
	if scheme == "" {
		if r.TLS == nil {
			scheme = "http"
		} else {
			scheme = "https"
		}
	}
	jwtStr, err := e.jwtGenerate(key, settings.ExplTokenValidity)
	if err != nil {
		return nil, err
	}
	return url.Parse(scheme + "://" + r.Host + e.webExplPathPrefix + jwtStr)
}
