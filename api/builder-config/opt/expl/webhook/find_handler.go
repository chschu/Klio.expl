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

func NewFindHandler(edb expldb.Finder, webFindPathPrefix string, jwtGenerator security.JwtGenerator, entryStringer types.EntryStringer, maxImmediateResults int, webUrlValidity time.Duration) Handler {
	return &findHandler{
		edb:                 edb,
		webFindPathPrefix:   webFindPathPrefix,
		jwtGenerator:        jwtGenerator,
		entryStringer:       entryStringer,
		maxImmediateResults: maxImmediateResults,
		webUrlValidity:      webUrlValidity,
	}
}

type findHandler struct {
	edb               expldb.Finder
	webFindPathPrefix string
	jwtGenerator      security.JwtGenerator
	entryStringer     types.EntryStringer

	maxImmediateResults int
	webUrlValidity      time.Duration
}

func (f *findHandler) Handle(in *Request, r *http.Request, now time.Time) (*Response, error) {
	syntaxResponse := NewResponse(fmt.Sprintf("Syntax: %s <POSIX-Regex>", in.TriggerWord))

	sep := regexp.MustCompile("^\\pZ*\\PZ+\\pZ+(?P<Regex>\\PZ+)\\pZ*$")
	match := sep.FindStringSubmatch(in.Text)
	if match == nil {
		return syntaxResponse, nil
	}
	rex := match[sep.SubexpIndex("Regex")]

	entries, total, err := f.edb.FindWithLimit(r.Context(), rex, f.maxImmediateResults)
	if err == expldb.ErrFindRegexInvalid {
		return syntaxResponse, nil
	}
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
			findUrl, err := f.getWebFindUrl(r, rex, now)
			if err != nil {
				logrus.Warnf("unable to resolve URL for web find: %v", err)
				totalText = fmt.Sprintf("%d Einträge", total)
			} else {
				totalText = fmt.Sprintf("[%d Einträge](%s)", total, findUrl)
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
			sb.WriteString(f.entryStringer.String(&entry))
			sb.WriteRune('\n')
		}
		sb.WriteString("```")
	}

	return NewResponse(sb.String()), nil
}

func (f *findHandler) getWebFindUrl(r *http.Request, rex string, now time.Time) (*url.URL, error) {
	scheme := r.URL.Scheme
	if scheme == "" {
		if r.TLS == nil {
			scheme = "http"
		} else {
			scheme = "https"
		}
	}
	jwtStr, err := f.jwtGenerator.Generate(rex, now.Add(f.webUrlValidity))
	if err != nil {
		return nil, err
	}
	return url.Parse(scheme + "://" + r.Host + f.webFindPathPrefix + jwtStr)
}
