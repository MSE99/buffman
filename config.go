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
	loginInterval  time.Duration
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

	parsedPollIntr, pollErr := time.ParseDuration(getEnv("POLL_INTERVAL", "1s"))
	if pollErr != nil {
		log.Panic(pollErr)
	}
	pollInterval = parsedPollIntr

	loginRefreshIntr, loginErr := time.ParseDuration(getEnv("LOGIN_INTERVAL", "30m"))
	if loginErr != nil {
		log.Panic(loginErr)
	}
	loginInterval = loginRefreshIntr
}

func getEnv(key string, def ...string) string {
	val, found := os.LookupEnv(key)
	if !found && len(def) > 0 {
		return def[0]
	}
	return val
}
