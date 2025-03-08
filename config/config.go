package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	HttpPort        string
	FmaUsername     string
	FmaPassword     string
	FmaLoginURL     string
	FmaDispatchURL  string
	DbFile          string
	Env             string
	OdooSecret      string
	PollInterval    time.Duration
	LoginInterval   time.Duration
	ContinueOnError bool
)

func loadConfigFromEnv() {
	if os.Getenv("ENV") == "" || os.Getenv("ENV") == "dev" {
		godotenv.Load()
	}

	Env = getEnv("ENV", "dev")
	HttpPort = getEnv("PORT", "3000")
	FmaUsername = getEnv("FMA_USERNAME")
	FmaPassword = getEnv("FMA_PASSWORD")
	FmaLoginURL = getEnv("FMA_LOGIN_URL")
	FmaDispatchURL = getEnv("FMA_DISPATCH_URL")
	DbFile = getEnv("DB")
	OdooSecret = getEnv("ODOO_SECRET")
	ContinueOnError = getEnv("DISPATCH_STRATEGY", "break") == "continue"

	parsedPollIntr, pollErr := time.ParseDuration(getEnv("POLL_INTERVAL", "1s"))
	if pollErr != nil {
		log.Panic(pollErr)
	}
	PollInterval = parsedPollIntr

	loginRefreshIntr, loginErr := time.ParseDuration(getEnv("LOGIN_INTERVAL", "30m"))
	if loginErr != nil {
		log.Panic(loginErr)
	}
	LoginInterval = loginRefreshIntr
}

func getEnv(key string, def ...string) string {
	val, found := os.LookupEnv(key)
	if !found && len(def) > 0 {
		return def[0]
	}
	return val
}

func Load() {
	loadConfigFromEnv()
}
