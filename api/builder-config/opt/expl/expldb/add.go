package expldb

import (
	"klio/expl/types"
	"time"
)

func (e *ExplDB) Add(key string, value string, createdBy string, createdAt time.Time) (entry *types.Entry, err error) {
	entry = &types.Entry{}
	err = e.db.Get(entry,
		"INSERT INTO entry(key, value, created_by, created_at) VALUES($1, $2, $3, $4) RETURNING *",
		key, value, createdBy, createdAt)
	if err != nil {
		return nil, err
	}
	return entry, nil
}
