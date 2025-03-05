package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	loadConfigFromEnv()

	log.SetOutput(os.Stdout)
	log.Println("running in ", env)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	db, dbErr := connectToDB(ctx, dbFile)
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
		err := app.Listen(":" + httpPort)

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
