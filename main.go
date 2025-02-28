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

	app := fiber.New()

	app.Get("/status", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).Send([]byte("OK"))
	})

	go app.Listen(":3000")

	<-ctx.Done()

	err := app.ShutdownWithTimeout(time.Second * 30)
	if err != nil {
		log.Fatal(err)
	}
}
