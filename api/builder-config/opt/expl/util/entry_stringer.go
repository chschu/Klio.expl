package util

import (
	"fmt"
	"klio/expl/types"
	"regexp"
	"strings"
	"time"
)

type EntryStringer interface {
	String(*types.Entry) string
}

type EntryStringerSettings interface {
	EntryToStringTimeFormat() string
	EntryToStringLocation() *time.Location
}

func NewEntryStringer(settings EntryStringerSettings) EntryStringer {
	return &entryStringer{
		settings: settings,
	}
}

type entryStringer struct {
	settings EntryStringerSettings
}

func (s *entryStringer) String(e *types.Entry) string {
	text := regexp.MustCompile("[[:space:]]").ReplaceAllString(e.Value, " ")

	var metadata []string
	if e.CreatedBy.Valid {
		createdBy := e.CreatedBy.String
		metadata = append(metadata, createdBy)
	}
	if e.CreatedAt.Valid {
		createdAt := e.CreatedAt.Time.In(s.settings.EntryToStringLocation()).Format(s.settings.EntryToStringTimeFormat())
		metadata = append(metadata, createdAt)
	}
	metadataText := ""
	if len(metadata) > 0 {
		metadataText = " (" + strings.Join(metadata, ", ") + ")"
	}

	return fmt.Sprintf("%s[%s/%s]: %s%s", e.Key, e.HeadIndex, e.PermanentIndex, text, metadataText)
}
