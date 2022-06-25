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
	"strings"
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

	var sb strings.Builder
	if count == 0 {
		sb.WriteString("Ich habe leider keinen Eintrag gefunden.")
	} else {
		var limitedEntries []types.Entry
		if count > settings.MaxExplCount {
			var entriesText string
			explUrl, err := e.getWebExplUrl(r, key)
			if err != nil {
				logrus.Warnf("unable to resolve URL for web expl: %v", err)
				entriesText = fmt.Sprintf("%d Einträge", count)
			} else {
				entriesText = fmt.Sprintf("[%d Einträge](%s)", count, explUrl)
			}
			sb.WriteString(fmt.Sprintf("Ich habe %s gefunden, das sind die letzten %d:\n",
				entriesText, settings.MaxExplCount))
			limitedEntries = entries[count-settings.MaxExplCount:]
		} else {
			if count == 1 {
				sb.WriteString("Ich habe den folgenden Eintrag gefunden:\n")
			} else {
				sb.WriteString(fmt.Sprintf("Ich habe die folgenden %d Einträge gefunden:\n", count))
			}
			limitedEntries = entries
		}
		sb.WriteString("```\n")
		for _, entry := range limitedEntries {
			sb.WriteString(entry.String())
			sb.WriteRune('\n')
		}
		sb.WriteString("```")
	}

	return NewResponse(sb.String()), nil
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
