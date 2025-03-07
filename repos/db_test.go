package repos

import (
	"context"
	"testing"
)

var ctx = context.Background()

func Test_ConnectToDB(t *testing.T) {
	t.Parallel()

	db, err := ConnectToDB(
		ctx,
		":memory:",
	)
	if err != nil {
		t.Error(err)
	}

	closeErr := db.Close()
	if closeErr != nil {
		t.Error(closeErr)
	}
}
