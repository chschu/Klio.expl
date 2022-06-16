package expldb

import (
	"klio/expl/types"
)

func (e *ExplDB) Expl(key string, indexRanges []types.IndexRange) (entries []types.Entry, err error) {
	irc, params := indexRangesSqlCondition(indexRanges)
	sql := "SELECT * FROM entry WHERE (" + irc + ") AND key_normalized = NORMALIZE(LOWER(?), NFC)"
	params = append(params, key)
	err = e.db.Select(&entries, e.db.Rebind(sql), params...)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
