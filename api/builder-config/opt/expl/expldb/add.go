package expldb

import (
	"klio/expl/types"
	"klio/expl/util"
	"time"
)

func (e *ExplDB) Add(key string, value string, createdBy string, createdAt time.Time) (entry *types.Entry, err error) {
	stmt, err := e.db.Prepare(
		"INSERT INTO entry(key, value, created_by, created_at) " +
			"VALUES($1, $2, $3, $4) " +
			"RETURNING key, value, created_by, created_at, head_index, tail_index, permanent_index")
	if err != nil {
		return nil, err
	}
	defer util.CloseAndAppendError(stmt, &err)

	row := stmt.QueryRow(key, value, createdBy, createdAt)
	entry = &types.Entry{}
	err = row.Scan(&entry.Key, &entry.Value, &entry.CreatedBy, &entry.CreatedAt, &entry.HeadIndex, &entry.TailIndex, &entry.PermanentIndex)
	if err != nil {
		return nil, err
	}

	return entry, nil
}
