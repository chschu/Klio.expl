package webhook

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"klio/expl/types"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func NewEntryListStringer(jwtGenerator JWTGenerator, entryStringer EntryStringer) *entryListStringer {
	return &entryListStringer{
		jwtGenerator:  jwtGenerator,
		entryStringer: entryStringer,
	}
}

type entryListStringer struct {
	jwtGenerator  JWTGenerator
	entryStringer EntryStringer
}

func (e *entryListStringer) String(entries []types.Entry, total int, webUrlSubject string, webUrlPathPrefix string, webUrlExpiresAt time.Time, req *http.Request) string {
	count := len(entries)

	var sb strings.Builder
	if total == 0 {
		sb.WriteString("Ich habe leider keinen Eintrag gefunden.")
	} else {
		if total > count {
			var totalText string
			explUrl, err := e.getWebUrl(req, webUrlSubject, webUrlPathPrefix, webUrlExpiresAt)
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

	return sb.String()
}

func (e *entryListStringer) getWebUrl(r *http.Request, subject string, pathPrefix string, expiresAt time.Time) (*url.URL, error) {
	scheme := r.URL.Scheme
	if scheme == "" {
		if r.TLS == nil {
			scheme = "http"
		} else {
			scheme = "https"
		}
	}
	jwtStr, err := e.jwtGenerator.GenerateJWT(subject, expiresAt)
	if err != nil {
		return nil, err
	}
	return url.Parse(scheme + "://" + r.Host + pathPrefix + jwtStr)
}
