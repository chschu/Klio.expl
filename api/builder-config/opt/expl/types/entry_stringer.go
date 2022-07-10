package types

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func NewEntryStringer(timeFormat string, timeLocation *time.Location) *entryStringer {
	return &entryStringer{
		timeFormat:   timeFormat,
		timeLocation: timeLocation,
	}
}

type entryStringer struct {
	timeFormat   string
	timeLocation *time.Location
}

var entryStringerSpaceRegexp = regexp.MustCompile("[[:space:]]")

func (s *entryStringer) String(e *Entry) string {
	text := entryStringerSpaceRegexp.ReplaceAllString(e.Value, " ")

	var metadata []string
	if e.CreatedBy.Valid {
		createdBy := e.CreatedBy.String
		metadata = append(metadata, createdBy)
	}
	if e.CreatedAt.Valid {
		createdAt := e.CreatedAt.Time.In(s.timeLocation).Format(s.timeFormat)
		metadata = append(metadata, createdAt)
	}
	metadataText := ""
	if len(metadata) > 0 {
		metadataText = " (" + strings.Join(metadata, ", ") + ")"
	}

	return fmt.Sprintf("%s[%s/%s]: %s%s", e.Key, NewHeadIndex(e.HeadIndex), NewPermanentIndex(e.PermanentIndex), text, metadataText)
}
