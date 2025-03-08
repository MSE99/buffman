package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mse99/buffman/repos"
)

var ctx = context.Background()

func TestStatusEndpoint(t *testing.T) {
	t.Parallel()

	db, err := repos.ConnectToDB(ctx, ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	server := CreateServer(context.Background(), db)
	defer server.Shutdown()

	req := httptest.NewRequest(http.MethodGet, "/status", nil)

	res, err := server.Test(req)

	if err != nil {
		t.Fatal(err)
	} else if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 but got %d", res.StatusCode)
	}
}
