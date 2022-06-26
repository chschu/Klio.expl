package expldb

import (
	"context"
	"errors"
	"klio/expl/types"
)

func (e *ExplDB) Expl(ctx context.Context, key string, indexSpec types.IndexSpec) (entries []types.Entry, err error) {
	irc, params := indexSpecSqlCondition(indexSpec)
	sql := "SELECT * FROM entry WHERE (" + irc + ") AND key_normalized = NORMALIZE(LOWER(?), NFC) ORDER BY id"
	params = append(params, key)
	err = e.db.SelectContext(ctx, &entries, e.db.Rebind(sql), params...)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (e *ExplDB) ExplWithLimit(ctx context.Context, key string, indexSpec types.IndexSpec, limit int) (entries []types.Entry, total int, err error) {
	if limit <= 0 {
		// query cannot report total if limit is zero
		return nil, 0, errors.New("limit must be greater than zero")
	}

	irc, params := indexSpecSqlCondition(indexSpec)
	sql := "SELECT * FROM (SELECT *, COUNT(*) OVER() total FROM entry WHERE (" + irc +
		") AND key_normalized = NORMALIZE(LOWER(?), NFC) ORDER BY id DESC LIMIT ?) x ORDER BY id"
	params = append(params, key, limit)

	var extEntries []struct {
		types.Entry
		Total int `db:"total"`
	}
	err = e.db.SelectContext(ctx, &extEntries, e.db.Rebind(sql), params...)
	if err != nil {
		return nil, 0, err
	}

	for _, extEntry := range extEntries {
		entries = append(entries, extEntry.Entry)
		total = extEntry.Total
	}

	return entries, total, nil
}
