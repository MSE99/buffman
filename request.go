package main

import (
	"context"
	"database/sql"
	"time"
)

type request struct {
	Id        int       `json:"id"`
	Payload   string    `json:"string"`
	CreatedOn time.Time `json:"createdOn"`
}

func deleteRequestByID(ctx context.Context, db *sql.DB, id int) error {
	_, err := db.ExecContext(ctx, `DELETE FROM RequestsBacklog WHERE id = @id`, sql.Named("id", id))
	return err
}

func insertRequest(ctx context.Context, db *sql.DB, req request) (request, error) {
	row := db.QueryRowContext(
		ctx,
		`INSERT INTO RequestsBacklog (payload, createdOn) VALUES (@payload, @createdOn) RETURNING id, payload, createdOn`,
		sql.Named("payload", req.Payload),
		sql.Named("createdOn", req.CreatedOn),
	)

	scanErr := row.Scan(
		&req.Id,
		&req.Payload,
		&req.CreatedOn,
	)

	return req, scanErr
}

func loadUnfinishedRequests(ctx context.Context, db *sql.DB) ([]request, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, payload, createdOn FROM RequestsBacklog ORDER BY createdOn DESC`)
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
