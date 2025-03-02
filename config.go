package main

import (
	"os"

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
)

func loadConfigFromEnv() {
	if os.Getenv("env") == "dev" {
		godotenv.Load()
	}

	env = getEnv("ENV", "dev")
	httpPort = getEnv("PORT", "3000")
	fmaUsername = getEnv("FMA_USERNAME")
	fmaUsername = getEnv("FMA_PASSWORD")
	fmaUsername = getEnv("FMA_LOGIN_URL")
	fmaUsername = getEnv("FMA_DISPATCH_URL")
	fmaUsername = getEnv("DB")
}

func getEnv(key string, def ...string) string {
	val, found := os.LookupEnv(key)
	if !found && len(def) > 0 {
		return def[0]
	}
	return val
}
