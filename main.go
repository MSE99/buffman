package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/mse99/buffman/buffman"
	"github.com/mse99/buffman/config"
	"github.com/mse99/buffman/repos"
	"github.com/mse99/buffman/web"
)

func main() {
	config.Load()

	log.SetOutput(os.Stdout)
	log.Println("running in ", config.Env)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	db, dbErr := repos.ConnectToDB(ctx, config.DbFile)
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	defer db.Close()

	dispatchErr := buffman.StartDispatchToFMA(ctx, db)
	if dispatchErr != nil {
		log.Panic(dispatchErr)
	}

	app := web.CreateServer(ctx, db)

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
