package expldb

import (
	"errors"
	"klio/expl/types"
)

var ErrFindRegexInvalid = errors.New("not a valid POSIX regular expression")

func (e *ExplDB) Find(rex string) (entries []types.Entry, err error) {
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
	err = e.db.Select(&entries, e.db.Rebind(sql), rex)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
