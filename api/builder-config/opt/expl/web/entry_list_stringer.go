package web

import (
	"fmt"
	"github.com/chschu/Klio.expl/types"
	"strings"
)

func NewEntryListStringer(entryStringer EntryStringer) *entryListStringer {
	return &entryListStringer{
		entryStringer: entryStringer,
	}
}

type entryListStringer struct {
	entryStringer EntryStringer
}

func (e *entryListStringer) String(entries []types.Entry) string {
	count := len(entries)

	var sb strings.Builder
	if count == 0 {
		sb.WriteString("Ich habe leider keinen Eintrag gefunden.\n")
	} else {
		if count == 1 {
			sb.WriteString("Ich habe den folgenden Eintrag gefunden:\n")
		} else {
			sb.WriteString(fmt.Sprintf("Ich habe die folgenden %d Einträge gefunden:\n", count))
		}
		for _, entry := range entries {
			sb.WriteString(e.entryStringer.String(&entry))
			sb.WriteRune('\n')
		}
	}

	return sb.String()
}
