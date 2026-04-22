package service

import (
	"testing"

	"life-financial-assistant-backend/internal/model"
	"life-financial-assistant-backend/internal/repository"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupAuthServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("migrate schema: %v", err)
	}

	return db
}

func TestRefreshTokenSucceedsAfterLogin(t *testing.T) {
	db := setupAuthServiceTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() { repository.DB = originalDB }()

	auth := NewAuthService(&repository.UserRepository{}, "test-secret")
	if err := auth.Register("refresh-user", "password123", "refresh-user@example.com"); err != nil {
		t.Fatalf("register user: %v", err)
	}

	_, refreshToken, err := auth.Login("refresh-user", "password123")
	if err != nil {
		t.Fatalf("login user: %v", err)
	}

	_, _, err = auth.RefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("expected refresh token to succeed after login, got %v", err)
	}
}

func TestRefreshTokenFailsAfterLogout(t *testing.T) {
	db := setupAuthServiceTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() { repository.DB = originalDB }()

	auth := NewAuthService(&repository.UserRepository{}, "test-secret")
	if err := auth.Register("logout-user", "password123", "logout-user@example.com"); err != nil {
		t.Fatalf("register user: %v", err)
	}

	_, refreshToken, err := auth.Login("logout-user", "password123")
	if err != nil {
		t.Fatalf("login user: %v", err)
	}

	user, err := (&repository.UserRepository{}).GetByUsername("logout-user")
	if err != nil {
		t.Fatalf("get user: %v", err)
	}

	if err := auth.Logout(user.ID); err != nil {
		t.Fatalf("logout user: %v", err)
	}

	_, _, err = auth.RefreshToken(refreshToken)
	if err == nil {
		t.Fatal("expected refresh token to fail after logout")
	}
}
