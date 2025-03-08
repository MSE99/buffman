package web

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mse99/buffman/config"
	"github.com/mse99/buffman/repos"
)

var ctx = context.Background()

func createTestingServer(t *testing.T) (*fiber.App, *sql.DB) {
	db, err := repos.ConnectToDB(ctx, ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	server := CreateServer(context.Background(), db)
	t.Cleanup(func() { server.Shutdown() })

	return server, db
}

func TestStatusEndpoint(t *testing.T) {
	t.Parallel()

	server, _ := createTestingServer(t)

	req := httptest.NewRequest(http.MethodGet, "/status", nil)

	res, err := server.Test(req)

	if err != nil {
		t.Fatal(err)
	} else if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 but got %d", res.StatusCode)
	}
}

func TestHandleQueueRequest(t *testing.T) {
	config.OdooSecret = "HelloWorld"

	t.Run("InvalidSecret", func(t *testing.T) {
		server, _ := createTestingServer(t)

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res, resErr := server.Test(req)

		if resErr != nil {
			t.Error(resErr)
		} else if res.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected status 401 but got %d", res.StatusCode)
		}
	})

	t.Run("ValidSecretButInvalidBody", func(t *testing.T) {
		server, _ := createTestingServer(t)

		path := fmt.Sprintf("/?token=%s", config.OdooSecret)

		req := httptest.NewRequest(http.MethodPost, path, nil)
		res, resErr := server.Test(req)

		if resErr != nil {
			t.Error(resErr)
		} else if res.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status 400 but got %d", res.StatusCode)
		}
	})

	t.Run("HappyPath", func(t *testing.T) {
		server, _ := createTestingServer(t)

		path := fmt.Sprintf("/?token=%s", config.OdooSecret)

		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader("HelloWorld"))
		res, resErr := server.Test(req)

		if resErr != nil {
			t.Error(resErr)
		} else if res.StatusCode != http.StatusOK {
			t.Errorf("expected status 200 but got %d", res.StatusCode)
		}
	})
}
