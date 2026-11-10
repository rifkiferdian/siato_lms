package config

import (
	"strings"
	"testing"
)

func TestLoadRejectsShortSessionSecret(t *testing.T) {
	t.Setenv("SESSION_SECRET", "short")
	if _, err := Load(); err == nil {
		t.Fatal("expected short session secret to be rejected")
	}
}
func TestDatabaseDSNUsesUTC(t *testing.T) {
	d := Database{Host: "localhost", Port: "3306", Name: "lms", User: "root"}
	got := d.DSN()
	if got == "" || !strings.Contains(got, "loc=UTC") {
		t.Fatalf("unexpected DSN: %s", got)
	}
}
