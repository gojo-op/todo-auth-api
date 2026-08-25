package service

import (
	"testing"

	"github.com/gojo-op/todo-auth-api/internal/config"
	"github.com/gojo-op/todo-auth-api/internal/models"
	"github.com/gojo-op/todo-auth-api/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuthService(t *testing.T) *AuthService {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	users := repository.NewUserRepository(db)
	cfg := &config.Config{JWTSecret: "test-secret"}

	return NewAuthService(users, cfg)
}

func TestRegisterAndLogin(t *testing.T) {
	auth := setupAuthService(t)

	registerInput := RegisterInput{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "secret123",
	}

	registered, err := auth.Register(registerInput)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if registered.Token == "" {
		t.Fatal("expected token on register")
	}

	loggedIn, err := auth.Login(LoginInput{
		Email:    registerInput.Email,
		Password: registerInput.Password,
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if loggedIn.Token == "" {
		t.Fatal("expected token on login")
	}

	userID, err := auth.ValidateToken(loggedIn.Token)
	if err != nil {
		t.Fatalf("validate token failed: %v", err)
	}
	if userID != registered.User.ID {
		t.Fatalf("expected user id %d, got %d", registered.User.ID, userID)
	}
}

func TestLoginWithInvalidPassword(t *testing.T) {
	auth := setupAuthService(t)

	_, err := auth.Register(RegisterInput{
		Name:     "Test User",
		Email:    "wrong@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	_, err = auth.Login(LoginInput{
		Email:    "wrong@example.com",
		Password: "bad-password",
	})
	if err != ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
