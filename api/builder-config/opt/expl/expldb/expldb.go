package expldb

import (
	"github.com/jmoiron/sqlx"
)

type ExplDB struct {
	db *sqlx.DB
}
