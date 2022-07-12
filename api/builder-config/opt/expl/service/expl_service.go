package service

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"klio/expl/types"
	"time"
)

func NewExplService(db *sqlx.DB) ExplService {
	return ExplService{
		db: db,
	}
}

type ExplService struct {
	db *sqlx.DB
}

func (e ExplService) Close() error {
	return e.db.Close()
}

func (e ExplService) Add(ctx context.Context, key string, value string, createdBy string, createdAt time.Time) (entry *types.Entry, err error) {
	entry = &types.Entry{}
	sql := "INSERT INTO entry(key, value, created_by, created_at) VALUES(?, ?, ?, ?) RETURNING *"
	err = e.db.GetContext(ctx, entry, e.db.Rebind(sql), key, value, createdBy, createdAt)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (e ExplService) Explain(ctx context.Context, key string) (entries []types.Entry, err error) {
	sql := "SELECT * FROM entry WHERE key_normalized = NORMALIZE(LOWER(?), NFC) ORDER BY id"
	err = e.db.SelectContext(ctx, &entries, e.db.Rebind(sql), key)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (e ExplService) ExplainWithLimit(ctx context.Context, key string, indexSpec types.IndexSpec, limit int) (entries []types.Entry, total int, err error) {
	if limit <= 0 {
		// query cannot report total if limit is zero
		return nil, 0, errors.New("limit must be greater than zero")
	}

	irc, params := indexSpec.SQLCondition()
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

func (e ExplService) Delete(ctx context.Context, key string, indexSpec types.IndexSpec) (entries []types.Entry, err error) {
	irc, params := indexSpec.SQLCondition()
	sql := "DELETE FROM entry WHERE (" + irc + ") AND key_normalized = NORMALIZE(LOWER(?), NFC) RETURNING *"
	params = append(params, key)
	err = e.db.SelectContext(ctx, &entries, e.db.Rebind(sql), params...)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

var ErrFindRegexInvalid = errors.New("not a valid POSIX regular expression")

func (e ExplService) Find(ctx context.Context, rex string) (entries []types.Entry, err error) {
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

func (e ExplService) FindWithLimit(ctx context.Context, rex string, limit int) (entries []types.Entry, total int, err error) {
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

func (e ExplService) Top(ctx context.Context, count int) (entries []types.Entry, err error) {
	sql := "SELECT * FROM entry WHERE tail_index = 1 ORDER BY head_index DESC, key_normalized LIMIT ?"
	err = e.db.SelectContext(ctx, &entries, e.db.Rebind(sql), count)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
