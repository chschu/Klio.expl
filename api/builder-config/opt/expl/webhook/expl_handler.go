package webhook

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"klio/expl/expldb"
	"klio/expl/security"
	"klio/expl/types"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func NewExplHandler(edb expldb.Explainer, webExplPathPrefix string, indexSpecParser types.IndexSpecParser, jwtGenerator security.JwtGenerator, entryStringer types.EntryStringer, settings ExplHandlerSettings) Handler {
	return &explHandler{
		edb:               edb,
		webExplPathPrefix: webExplPathPrefix,
		jwtGenerator:      jwtGenerator,
		entryStringer:     entryStringer,
		settings:          settings,
		indexSpecParser:   indexSpecParser,
	}
}

type ExplHandlerSettings interface {
	MaxExplCount() int
	ExplTokenValidity() time.Duration
}

type explHandler struct {
	edb               expldb.Explainer
	webExplPathPrefix string
	indexSpecParser   types.IndexSpecParser
	jwtGenerator      security.JwtGenerator
	entryStringer     types.EntryStringer
	settings          ExplHandlerSettings
}

func (e *explHandler) Handle(in *Request, r *http.Request, now time.Time) (*Response, error) {
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
		indexSpec, err = e.indexSpecParser.ParseIndexSpec(indexSpecStr)
		if err != nil {
			return syntaxResponse, nil
		}
	}

	entries, total, err := e.edb.ExplainWithLimit(r.Context(), key, indexSpec, e.settings.MaxExplCount())
	if err != nil {
		return nil, err
	}

	count := len(entries)

	var sb strings.Builder
	if total == 0 {
		sb.WriteString("Ich habe leider keinen Eintrag gefunden.")
	} else {
		if total > count {
			var totalText string
			explUrl, err := e.getWebExplUrl(r, key, now)
			if err != nil {
				logrus.Warnf("unable to resolve URL for web expl: %v", err)
				totalText = fmt.Sprintf("%d Einträge", total)
			} else {
				totalText = fmt.Sprintf("[%d Einträge](%s)", total, explUrl)
			}
			sb.WriteString(fmt.Sprintf("Ich habe %s gefunden, das sind die letzten %d:\n", totalText, count))
		} else {
			if total == 1 {
				sb.WriteString("Ich habe den folgenden Eintrag gefunden:\n")
			} else {
				sb.WriteString(fmt.Sprintf("Ich habe die folgenden %d Einträge gefunden:\n", total))
			}
		}
		sb.WriteString("```\n")
		for _, entry := range entries {
			sb.WriteString(e.entryStringer.String(&entry))
			sb.WriteRune('\n')
		}
		sb.WriteString("```")
	}

	return NewResponse(sb.String()), nil
}

func (e *explHandler) getWebExplUrl(r *http.Request, key string, now time.Time) (*url.URL, error) {
	scheme := r.URL.Scheme
	if scheme == "" {
		if r.TLS == nil {
			scheme = "http"
		} else {
			scheme = "https"
		}
	}
	jwtStr, err := e.jwtGenerator.Generate(key, now.Add(e.settings.ExplTokenValidity()))
	if err != nil {
		return nil, err
	}
	return url.Parse(scheme + "://" + r.Host + e.webExplPathPrefix + jwtStr)
}
