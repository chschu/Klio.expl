package expldb

import (
	"context"
	"errors"
	"klio/expl/types"
)

var ErrFindRegexInvalid = errors.New("not a valid POSIX regular expression")

func (e *ExplDB) Find(ctx context.Context, rex string) (entries []types.Entry, err error) {
	var valid bool
	sql := "SELECT regexp_is_valid(?)"
	err = e.db.Get(&valid, e.db.Rebind(sql), rex)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, ErrFindRegexInvalid
	}
	sql = "SELECT * FROM entry WHERE NORMALIZE(value, NFC) ~ NORMALIZE(?, NFC) ORDER BY id"
	err = e.db.SelectContext(ctx, &entries, e.db.Rebind(sql), rex)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (e *ExplDB) FindWithLimit(ctx context.Context, rex string, limit int) (entries []types.Entry, total int, err error) {
	if limit <= 0 {
		// query cannot report total if limit is zero
		return nil, 0, errors.New("limit must be greater than zero")
	}

	var valid bool
	sql := "SELECT regexp_is_valid(?)"
	err = e.db.Get(&valid, e.db.Rebind(sql), rex)
	if err != nil {
		return nil, 0, err
	}
	if !valid {
		return nil, 0, ErrFindRegexInvalid
	}

	sql = "SELECT * FROM (SELECT *, COUNT(*) OVER() total FROM entry WHERE NORMALIZE(value, NFC) ~ NORMALIZE(?, NFC) ORDER BY id DESC LIMIT ?) x ORDER BY id"

	var extEntries []struct {
		types.Entry
		Total int `db:"total"`
	}
	err = e.db.SelectContext(ctx, &extEntries, e.db.Rebind(sql), rex, limit)
	if err != nil {
		return nil, 0, err
	}

	for _, extEntry := range extEntries {
		entries = append(entries, extEntry.Entry)
		total = extEntry.Total
	}

	return entries, total, nil
}
