package importer

import (
	"database/sql"
	_ "embed"
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"io/fs"
	"klio/expl/util"
	_ "modernc.org/sqlite"
	"os"
	"strings"
	"time"
)

type ImportEntry struct {
	Key       string  `db:"key"`
	Value     string  `db:"value"`
	CreatedBy *string `db:"createdBy"`
	CreatedAt *int64  `db:"createdAt"`
	Visible   uint8   `db:"visible"`
}

func Import() (err error) {
	_, err = os.Stat("/expl.sqlite")
	if errors.Is(err, fs.ErrNotExist) {
		logrus.Info("Import skipped because /expl.sqlite does not exist")
		return nil
	}
	if err != nil {
		return err
	}

	fromDb, err := sqlx.Open("sqlite", "file:///expl.sqlite")
	if err != nil {
		return err
	}
	defer util.CloseAndAppendError(fromDb, &err)

	toDb, err := sqlx.Open("postgres", os.Getenv("CONNECT_STRING"))
	if err != nil {
		return err
	}
	defer util.CloseAndAppendError(toDb, &err)

	empty, err := checkEmpty(toDb)
	if err != nil {
		return err
	}
	if !empty {
		logrus.Info("Import skipped because target database is not empty")
		return nil
	}

	entries, err := selectEntries(fromDb)
	if err != nil {
		return err
	}

	err = insertEntries(toDb, entries)
	if err != nil {
		return err
	}
	logrus.Infof("%d entries successfully imported", len(entries))

	return nil
}

func checkEmpty(toDb *sqlx.DB) (empty bool, err error) {
	q := "SELECT 1 FROM entry_data LIMIT 1"
	res, err := toDb.Query(q)
	defer util.CloseAndAppendError(res, &err)

	return !res.Next(), nil
}

func selectEntries(fromDb *sqlx.DB) (entries []ImportEntry, err error) {
	q := "SELECT item key, expl value, nick createdBy, datetime-62167219200000 createdAt, enabled visible " +
		"FROM t_expl ORDER BY id"
	err = fromDb.Select(&entries, q)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func insertEntries(toDb *sqlx.DB, entries []ImportEntry) error {
	var paramSql []string
	var paramValues []any

	for _, entry := range entries {
		paramSql = append(paramSql, "(?, NORMALIZE(LOWER(?), NFC), ?, ?, ?, ?)")
		paramValues = append(paramValues, entry.Key, entry.Key, entry.Value)
		if entry.CreatedBy != nil {
			paramValues = append(paramValues, sql.NullString{String: *entry.CreatedBy, Valid: true})
		} else {
			paramValues = append(paramValues, sql.NullString{})
		}
		if entry.CreatedAt != nil {
			paramValues = append(paramValues, sql.NullTime{Time: time.UnixMilli(*entry.CreatedAt), Valid: true})
		} else {
			paramValues = append(paramValues, sql.NullTime{})
		}
		paramValues = append(paramValues, entry.Visible != 0)

		if len(paramValues) >= 32768 {
			err := bulkInsert(toDb, paramSql, paramValues)
			if err != nil {
				return err
			}
			paramSql = paramSql[:0]
			paramValues = paramValues[:0]
		}
	}

	return bulkInsert(toDb, paramSql, paramValues)
}

func bulkInsert(toDb *sqlx.DB, paramSql []string, paramValues []any) error {
	if len(paramSql) == 0 {
		return nil
	}
	ins := "INSERT INTO entry_data(key, key_normalized, value, created_by, created_at, visible) " +
		"VALUES " + strings.Join(paramSql, ", ")
	_, err := toDb.Exec(toDb.Rebind(ins), paramValues...)
	return err
}
