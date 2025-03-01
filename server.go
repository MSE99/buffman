package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func createHttpServer(ctx context.Context, db *sql.DB) *fiber.App {
	app := fiber.New()
	setupRouter(ctx, app, db)
	return app
}

func setupRouter(ctx context.Context, app *fiber.App, db *sql.DB) {
	app.Get("/status", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).Send([]byte("OK"))
	})

	app.Post("/", func(c *fiber.Ctx) error {
		payload := string(c.Body())

		_, err := insertRequest(ctx, db, request{
			Payload:   payload,
			CreatedOn: time.Now(),
		})
		if err != nil {
			log.Println(err)
			return c.Status(http.StatusInternalServerError).Send([]byte(""))
		}

		processRequestsNow <- true

		return c.Status(http.StatusOK).Send([]byte("OK"))
	})
}
