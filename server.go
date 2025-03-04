package main

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func createHttpServer(ctx context.Context, db *sql.DB) *fiber.App {
	app := fiber.New()
	setupRouter(ctx, app, db)
	return app
}

func setupRouter(ctx context.Context, app *fiber.App, db *sql.DB) {
	app.Use(logger.New(logger.Config{Output: os.Stdout}))

	app.Get("/status", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).Send([]byte("OK"))
	})

	app.Post("/", func(c *fiber.Ctx) error {
		token := c.Params("token")
		hashedToken := sha256.Sum256([]byte(token))

		if subtle.ConstantTimeCompare(hashedToken[:], []byte(odooSecret)) == 0 {
			log.Println("received request with invalid secret")
			return c.Status(http.StatusUnauthorized).Send([]byte("Unauthorized"))
		}

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
