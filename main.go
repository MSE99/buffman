package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

var dispatchChan = make(chan request)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config, configErr := loadConfigFile()
	if configErr != nil {
		log.Fatal(configErr)
	}

	db, dbErr := connectToDB(ctx, config.DB)
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	defer db.Close()

	go listenForDispatches(ctx, db, dispatchChan)

	app := createHttpServer()
	go app.Listen(config.getHttpServerAddr())

	<-ctx.Done()

	err := app.ShutdownWithTimeout(time.Second * 30)
	if err != nil {
		log.Fatal(err)
	}
}
