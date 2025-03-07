package web

import (
	"context"
	"database/sql"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func CreateServer(ctx context.Context, db *sql.DB) *fiber.App {
	app := fiber.New()
	setupRouter(ctx, app, db)
	return app
}

func setupRouter(ctx context.Context, app *fiber.App, db *sql.DB) {
	app.Use(logger.New(logger.Config{Output: os.Stdout}))

	app.Get("/status", handleGetStatusRequest)
	app.Post("/", createQueueRequestHandler(ctx, db))
}
