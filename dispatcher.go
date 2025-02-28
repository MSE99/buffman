package main

import (
	"context"
	"database/sql"
	"log"
)

// TODO
func dispatchRequest(_ context.Context, _ *sql.DB, _ request) error {
	return nil
}

func listenForDispatches(ctx context.Context, db *sql.DB, requests <-chan request) {
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-requests:
			if req.Id == 0 {
				var err error
				req, err = insertRequest(ctx, db, req)

				if err != nil {
					log.Println("error while dispatching", err)
					continue
				}
			}

			dispatchErr := dispatchRequest(ctx, db, req)

			if dispatchErr != nil {
				log.Println("error while dispatching", dispatchErr)
				continue
			}

			_ = deleteRequestByID(ctx, db, req.Id)
		}
	}
}
