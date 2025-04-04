package repos

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectToDB(ctx context.Context, dbFilename string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		return nil, err
	}

	_, execErr := conn.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS RequestsBacklog (
			id INTEGER PRIMARY KEY,
			payload TEXT,
			createdOn DATETIME
		);
	`)
	if execErr != nil {
		return nil, execErr
	}

	return conn, nil
}
