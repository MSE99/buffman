package buffman

import (
	"context"
	"database/sql"
)

func StartDispatchToFMA(ctx context.Context, db *sql.DB) error {
	tk, err := newFmaToken(ctx)
	if err != nil {
		return err
	}
	go tk.waitAndRefresh()

	go processStoredRequests(ctx, requestProcessingOpts{
		db: db,
		tk: tk,
	})

	return nil
}
