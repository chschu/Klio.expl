package types

import (
	"time"
)

type Entry struct {
	Id             int32          `db:"id"`
	Key            string         `db:"key"`
	KeyNormalized  string         `db:"key_normalized"`
	Value          string         `db:"value"`
	CreatedBy      string         `db:"created_by"`
	CreatedAt      time.Time      `db:"created_at"`
	HeadIndex      HeadIndex      `db:"head_index"`
	TailIndex      TailIndex      `db:"tail_index"`
	PermanentIndex PermanentIndex `db:"permanent_index"`
}
