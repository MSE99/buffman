package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	t.Run("NoDefault", func(t *testing.T) {
		os.Setenv("foo", "baz")

		val := getEnv("foo")

		if val != "baz" {
			t.Errorf("expected %s to equal %s", val, "baz")
		}
	})

	t.Run("WithDefault", func(t *testing.T) {
		val := getEnv("foo", "baz")

		if val != "baz" {
			t.Errorf("expected %s to equal %s", val, "baz")
		}
	})
}

func TestLoadConfigFromEnv(t *testing.T) {
	prevEnv := getEnv("ENV", "dev")
	prevHttpPort := getEnv("PORT", "3000")
	prevFmaUsername := getEnv("FMA_USERNAME")
	prevFmaPassword := getEnv("FMA_PASSWORD")
	prevFmaLoginURL := getEnv("FMA_LOGIN_URL")
	prevFmaDispatchURL := getEnv("FMA_DISPATCH_URL")
	prevDbFile := getEnv("DB")
	prevOdooSecret := getEnv("ODOO_SECRET")

	t.Cleanup(func() {
		os.Setenv("ENV", prevEnv)
		os.Setenv("PORT", prevHttpPort)
		os.Setenv("FMA_USERNAME", prevFmaUsername)
		os.Setenv("FMA_PASSWORD", prevFmaPassword)
		os.Setenv("FMA_LOGIN_URL", prevFmaLoginURL)
		os.Setenv("FMA_DISPATCH_URL", prevFmaDispatchURL)
		os.Setenv("DB", prevDbFile)
		os.Setenv("ODOO_SECRET", prevOdooSecret)
	})

	t.Run("LoadingFromDotEnv", func(t *testing.T) {
		loadConfigFromEnv()
	})

	t.Run("AllSet", func(t *testing.T) {
		os.Setenv("ENV", "prod")
		os.Setenv("PORT", "3500")
		os.Setenv("FMA_USERNAME", "admin")
		os.Setenv("FMA_PASSWORD", "admin")
		os.Setenv("FMA_LOGIN_URL", "login")
		os.Setenv("FMA_DISPATCH_URL", "dispatch")
		os.Setenv("DB", "FILO.db")
		os.Setenv("ODOO_SECRET", "FOO")

		loadConfigFromEnv()

		if HttpPort != "3500" {
			t.Errorf("expected httpPort to be 3500 but got, %s", HttpPort)
		}

		if FmaUsername != "admin" {
			t.Errorf("expected fmaUsername to be admin but got, %s", FmaUsername)
		}

		if FmaPassword != "admin" {
			t.Errorf("expected fmaPassword to be admin but got, %s", FmaPassword)
		}

		if FmaLoginURL != "login" {
			t.Errorf("expected fmaLoginURL to be login but got, %s", FmaLoginURL)
		}

		if FmaDispatchURL != "dispatch" {
			t.Errorf("expected fmaDispatchURL to be dispatch but got, %s", FmaDispatchURL)
		}

		if DbFile != "FILO.db" {
			t.Errorf("expected dbFile to be FILO.db but got, %s", DbFile)
		}

		if OdooSecret != "FOO" {
			t.Errorf("expected odooSecret to be FOO but got, %s", OdooSecret)
		}
	})
}
