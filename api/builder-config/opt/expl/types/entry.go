package types

import (
	"time"
)

type Entry struct {
	Key            string
	Value          string
	CreatedBy      string
	CreatedAt      time.Time
	HeadIndex      HeadIndex
	TailIndex      TailIndex
	PermanentIndex PermanentIndex
}
