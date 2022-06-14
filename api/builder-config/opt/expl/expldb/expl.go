package expldb

import (
	"klio/expl/types"
)

func (e *ExplDB) Expl(key string, indexRanges []types.IndexRange) (entries []types.Entry, err error) {
	irc, params := indexRangesSqlCondition(indexRanges, []any{key})
	err = e.db.Select(&entries,
		"SELECT * FROM entry WHERE key_normalized = NORMALIZE(LOWER($1), NFC) AND ("+irc+")",
		params...)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
