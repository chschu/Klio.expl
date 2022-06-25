package expldb

import (
	"klio/expl/types"
)

func (e *ExplDB) Expl(key string, indexSpec types.IndexSpec) (entries []types.Entry, err error) {
	irc, params := indexSpecSqlCondition(indexSpec)
	sql := "SELECT * FROM entry WHERE (" + irc + ") AND key_normalized = NORMALIZE(LOWER(?), NFC) ORDER BY id"
	params = append(params, key)
	err = e.db.Select(&entries, e.db.Rebind(sql), params...)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
