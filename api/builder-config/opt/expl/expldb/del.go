package expldb

import (
	"klio/expl/types"
)

func (e *ExplDB) Del(key string, indexRanges []types.IndexRange) (entries []types.Entry, err error) {
	irc, params := indexRangesSqlCondition(indexRanges, []any{key})
	err = e.db.Select(&entries,
		"DELETE FROM entry WHERE key_normalized = NORMALIZE(LOWER($1), NFC) AND ("+irc+") RETURNING *",
		params...)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
