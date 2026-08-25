package config

import (
	"os"
	"testing"
)

func TestLoadRequiresSqlitePath(t *testing.T) {
	t.Setenv("SQLITE_PATH", "")
	t.Setenv("SERVER_PORT", "8080")
	t.Setenv("JWT_SECRET", "test-secret")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when SQLITE_PATH is missing")
	}
}

func TestLoadRequiresJWTSecret(t *testing.T) {
	t.Setenv("SQLITE_PATH", "test.db")
	t.Setenv("SERVER_PORT", "8080")
	t.Setenv("JWT_SECRET", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when JWT_SECRET is missing")
	}
}

func TestLoadDefaultsServerPort(t *testing.T) {
	t.Setenv("SQLITE_PATH", "test.db")
	t.Setenv("SERVER_PORT", "")
	t.Setenv("JWT_SECRET", "test-secret")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.ServerPort != "8080" {
		t.Fatalf("expected default port 8080, got %s", cfg.ServerPort)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
