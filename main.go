package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/mse99/buffman/config"
)

func main() {
	config.Load()

	log.SetOutput(os.Stdout)
	log.Println("running in ", config.Env)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	db, dbErr := connectToDB(ctx, config.DbFile)
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	defer db.Close()

	tk, tkErr := newFmaToken(ctx)
	if tkErr != nil {
		log.Fatal(tkErr)
	}
	go tk.waitAndRefresh()

	go processStoredRequests(ctx, requestProcessingOpts{
		db: db,
		tk: tk,
	})

	app := createHttpServer(ctx, db)

	go func() {
		err := app.Listen(":" + config.HttpPort)

		if err != nil {
			log.Println("failed to listen", err)
		}
	}()

	<-ctx.Done()

	err := app.ShutdownWithTimeout(time.Second * 30)
	if err != nil {
		log.Fatal(err)
	}
}
