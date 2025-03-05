package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func assertNotErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("gotten error %v", err)
	}
}

func TestStatus(t *testing.T) {
	db := connectToTestingDB(t)

	app := fiber.New()
	defer app.Shutdown()

	setupRouter(ctx, app, db)

	req := httptest.NewRequest(http.MethodGet, "/status", nil)

	res, err := app.Test(req)
	assertNotErr(t, err)

	if res.StatusCode != http.StatusOK {
		t.Errorf("server responded with non 200 status %d", res.StatusCode)
	}
}
