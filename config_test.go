package main

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	t.Parallel()

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
