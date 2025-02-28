package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	db, dbErr := connectToDB(ctx, "buffman.db")
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	defer db.Close()

	dispatchChan := make(chan request)

	go listenForDispatches(ctx, db, dispatchChan)

	app := fiber.New()

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

	go app.Listen(":3000")

	<-ctx.Done()

	err := app.ShutdownWithTimeout(time.Second * 30)
	if err != nil {
		log.Fatal(err)
	}
}
