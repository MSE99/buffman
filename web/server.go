package web

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
	"github.com/mse99/buffman/buffman"
	"github.com/mse99/buffman/config"
)

func CreateServer(ctx context.Context, db *sql.DB) *fiber.App {
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
		token := c.Query("token")
		hashedToken := sha256.Sum256([]byte(token))
		hashedOdooSecret := sha256.Sum256([]byte(config.OdooSecret))

		if subtle.ConstantTimeCompare(hashedToken[:], hashedOdooSecret[:]) == 0 {
			log.Println("received request with invalid secret")
			return c.Status(http.StatusUnauthorized).Send([]byte("Unauthorized"))
		}

		payload := string(c.Body())

		if len(payload) == 0 {
			return c.Status(http.StatusBadRequest).Send([]byte("Invalid body sent"))
		}

		queueErr := buffman.QueueRequest(ctx, db, buffman.Request{Payload: payload, CreatedOn: time.Now()})
		if queueErr != nil {
			log.Println(queueErr)
			return c.Status(http.StatusInternalServerError).Send([]byte(""))
		}

		return c.Status(http.StatusOK).Send([]byte("OK"))
	})

}
