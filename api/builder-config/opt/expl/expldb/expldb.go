package expldb

import (
	"database/sql"
)

type ExplDB struct {
	db *sql.DB
}
