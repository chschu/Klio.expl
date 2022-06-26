package expldb

import (
	"context"
	"klio/expl/types"
)

func (e *ExplDB) Top(ctx context.Context, count int) (entries []types.Entry, err error) {
	sql := "SELECT * FROM entry WHERE tail_index = 1 ORDER BY head_index DESC, key_normalized LIMIT ?"
	err = e.db.SelectContext(ctx, &entries, e.db.Rebind(sql), count)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
