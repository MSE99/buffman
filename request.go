package main

import (
	"context"
	"database/sql"
	"time"
)

type request struct {
	Id        int        `json:"id"`
	Payload   string     `json:"string"`
	CreatedOn *time.Time `json:"createdOn"`
}

func deleteRequestByID(ctx context.Context, db *sql.DB, id int) error {
	_, err := db.ExecContext(ctx, `DELETE FROM RequestsBacklog WHERE id = @id`, sql.Named("id", id))
	return err
}

func insertRequest(ctx context.Context, db *sql.DB, req request) error {
	_, err := db.ExecContext(
		ctx,
		`INSERT INTO RequestsBacklog (id, payload, createdOn) VALUES (@id, @@payload, @createdOn)`,
		sql.Named("id", req.Id),
		sql.Named("payload", req.Payload),
		sql.Named("createdOn", req.CreatedOn),
	)

	return err
}

func loadUnfinishedRequests(ctx context.Context, db *sql.DB) ([]request, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, payload, createdOn FROM requests ORDER BY createdOn DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []request{}

	for rows.Next() {
		req := request{}

		scanErr := rows.Scan(
			&req.Id,
			&req.Payload,
			&req.CreatedOn,
		)
		if scanErr != nil {
			return nil, scanErr
		}

		results = append(results, req)
	}

	return results, nil
}
