package types

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"klio/expl/settings"
	"regexp"
	"strings"
	"time"
)

type Entry struct {
	Id             int32          `db:"id"`
	Key            string         `db:"key"`
	KeyNormalized  string         `db:"key_normalized"`
	Value          string         `db:"value"`
	CreatedBy      sql.NullString `db:"created_by"`
	CreatedAt      sql.NullTime   `db:"created_at"`
	HeadIndex      HeadIndex      `db:"head_index"`
	TailIndex      TailIndex      `db:"tail_index"`
	PermanentIndex PermanentIndex `db:"permanent_index"`
}

func (e *Entry) String() string {
	text := regexp.MustCompile("[[:space:]]").ReplaceAllString(e.Value, " ")

	var metadata []string
	if e.CreatedBy.Valid {
		createdBy := e.CreatedBy.String
		metadata = append(metadata, createdBy)
	}
	if e.CreatedAt.Valid {
		createdAt := e.CreatedAt.Time.In(entryToStringTimeZone).Format(settings.EntryToStringTimeFormat)
		metadata = append(metadata, createdAt)
	}
	metadataText := ""
	if len(metadata) > 0 {
		metadataText = " (" + strings.Join(metadata, ", ") + ")"
	}

	return fmt.Sprintf("%s[%s/%s]: %s%s", e.Key, e.HeadIndex, e.PermanentIndex, text, metadataText)
}

var entryToStringTimeZone = func(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		logrus.Fatal(err)
	}
	return loc
}(settings.EntryToStringTimeZone)
