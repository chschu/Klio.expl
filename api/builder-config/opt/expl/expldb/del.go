package expldb

import (
	"context"
	"klio/expl/types"
)

func (e *ExplDB) Del(ctx context.Context, key string, indexSpec types.IndexSpec) (entries []types.Entry, err error) {
	irc, params := indexSpecSqlCondition(indexSpec)
	sql := "DELETE FROM entry WHERE (" + irc + ") AND key_normalized = NORMALIZE(LOWER(?), NFC) RETURNING *"
	params = append(params, key)
	err = e.db.SelectContext(ctx, &entries, e.db.Rebind(sql), params...)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
