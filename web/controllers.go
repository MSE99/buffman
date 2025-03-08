package web

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mse99/buffman/buffman"
	"github.com/mse99/buffman/config"
)

func handleGetStatusRequest(ctx *fiber.Ctx) error {
	return ctx.Status(200).Send([]byte("OK"))
}

func createQueueRequestHandler(ctx context.Context, db *sql.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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

		queueCtx, cancel := context.WithTimeout(ctx, time.Millisecond*50)
		defer cancel()

		queueErr := buffman.QueueRequest(queueCtx, db, payload)
		if queueErr != nil {
			log.Println("Error while attempting to queue request", queueErr)
			return c.Status(http.StatusInternalServerError).Send([]byte(""))
		}

		return c.Status(http.StatusOK).Send([]byte("OK"))
	}
}
