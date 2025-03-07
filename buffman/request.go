package buffman

import (
	"context"
	"database/sql"
	"time"
)

type Request struct {
	Id        int       `json:"id"`
	Payload   string    `json:"string"`
	CreatedOn time.Time `json:"createdOn"`
}

func deleteRequestByID(ctx context.Context, db *sql.DB, id int) error {
	_, err := db.ExecContext(ctx, `DELETE FROM RequestsBacklog WHERE id = @id`, sql.Named("id", id))
	return err
}

func insertRequest(ctx context.Context, db *sql.DB, req Request) (Request, error) {
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

func loadUnfinishedRequests(ctx context.Context, db *sql.DB) ([]Request, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, payload, createdOn FROM RequestsBacklog ORDER BY createdOn ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []Request{}

	for rows.Next() {
		req := Request{}

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
