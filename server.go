package main

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func createHttpServer() *fiber.App {
	app := fiber.New()
	setupRouter(app)
	return app
}

func setupRouter(app *fiber.App) {
	app.Get("/status", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).Send([]byte("OK"))
	})

	app.Post("/", func(c *fiber.Ctx) error {
		payload := string(c.Body())

		req := request{
			Payload:   payload,
			CreatedOn: time.Now(),
		}

		dispatchChan <- req

		return c.Status(http.StatusOK).Send([]byte("OK"))
	})
}
