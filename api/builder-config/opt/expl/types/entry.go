package types

import (
	"database/sql"
)

type Entry struct {
	Id             int32          `db:"id"`
	Key            string         `db:"key"`
	KeyNormalized  string         `db:"key_normalized"`
	Value          string         `db:"value"`
	CreatedBy      sql.NullString `db:"created_by"`
	CreatedAt      sql.NullTime   `db:"created_at"`
	HeadIndex      uint           `db:"head_index"`
	TailIndex      uint           `db:"tail_index"`
	PermanentIndex uint           `db:"permanent_index"`
}
