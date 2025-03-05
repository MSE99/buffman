package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	httpPort       string
	fmaUsername    string
	fmaPassword    string
	fmaLoginURL    string
	fmaDispatchURL string
	dbFile         string
	env            string
	odooSecret     string
	pollInterval   time.Duration
)

func loadConfigFromEnv() {
	if os.Getenv("ENV") == "" || os.Getenv("ENV") == "dev" {
		godotenv.Load()
	}

	env = getEnv("ENV", "dev")
	httpPort = getEnv("PORT", "3000")
	fmaUsername = getEnv("FMA_USERNAME")
	fmaPassword = getEnv("FMA_PASSWORD")
	fmaLoginURL = getEnv("FMA_LOGIN_URL")
	fmaDispatchURL = getEnv("FMA_DISPATCH_URL")
	dbFile = getEnv("DB")
	odooSecret = getEnv("ODOO_SECRET")

	parsedPollIntr, err := time.ParseDuration(getEnv("POLL_INTERVAL", "1s"))
	if err != nil {
		log.Panic(err)
	}
	pollInterval = parsedPollIntr
}

func getEnv(key string, def ...string) string {
	val, found := os.LookupEnv(key)
	if !found && len(def) > 0 {
		return def[0]
	}
	return val
}
