package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

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

	tk, tkErr := newFmaToken(ctx, &config)
	if tkErr != nil {
		log.Fatal(tkErr)
	}
	go tk.waitAndRefresh(&config)

	go processStoredRequests(ctx, requestProcessingOpts{
		db:           db,
		pollInterval: time.Second * 5,
	})

	app := createHttpServer(ctx, db)
	go app.Listen(config.getHttpServerAddr())

	<-ctx.Done()

	err := app.ShutdownWithTimeout(time.Second * 30)
	if err != nil {
		log.Fatal(err)
	}
}
