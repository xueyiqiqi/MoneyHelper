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

func TestUpdateProfileUpdatesUsernameAndEmail(t *testing.T) {
	db := setupAuthServiceTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() { repository.DB = originalDB }()

	auth := NewAuthService(&repository.UserRepository{}, "test-secret")
	if err := auth.Register("profile-user", "password123", "profile-user@example.com"); err != nil {
		t.Fatalf("register user: %v", err)
	}

	user, err := (&repository.UserRepository{}).GetByUsername("profile-user")
	if err != nil {
		t.Fatalf("get user: %v", err)
	}

	updated, err := auth.UpdateProfile(user.ID, "profile-user-2", "profile-user-2@example.com")
	if err != nil {
		t.Fatalf("update profile: %v", err)
	}

	if updated.Username != "profile-user-2" {
		t.Fatalf("expected updated username, got %s", updated.Username)
	}
	if updated.Email != "profile-user-2@example.com" {
		t.Fatalf("expected updated email, got %s", updated.Email)
	}
}

func TestUpdateProfileRejectsDuplicateEmail(t *testing.T) {
	db := setupAuthServiceTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() { repository.DB = originalDB }()

	auth := NewAuthService(&repository.UserRepository{}, "test-secret")
	if err := auth.Register("first-user", "password123", "first@example.com"); err != nil {
		t.Fatalf("register first user: %v", err)
	}
	if err := auth.Register("second-user", "password123", "second@example.com"); err != nil {
		t.Fatalf("register second user: %v", err)
	}

	second, err := (&repository.UserRepository{}).GetByUsername("second-user")
	if err != nil {
		t.Fatalf("get second user: %v", err)
	}

	_, err = auth.UpdateProfile(second.ID, "second-user", "first@example.com")
	if err == nil {
		t.Fatal("expected duplicate email update to fail")
	}
}

func TestChangePasswordRejectsWrongCurrentPassword(t *testing.T) {
	db := setupAuthServiceTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() { repository.DB = originalDB }()

	auth := NewAuthService(&repository.UserRepository{}, "test-secret")
	if err := auth.Register("pwd-user", "password123", "pwd-user@example.com"); err != nil {
		t.Fatalf("register user: %v", err)
	}

	user, err := (&repository.UserRepository{}).GetByUsername("pwd-user")
	if err != nil {
		t.Fatalf("get user: %v", err)
	}

	err = auth.ChangePassword(user.ID, "wrong-password", "new-password-123")
	if err == nil {
		t.Fatal("expected wrong current password to fail")
	}
}

func TestChangePasswordRevokesRefreshTokenAndAllowsLoginWithNewPassword(t *testing.T) {
	db := setupAuthServiceTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() { repository.DB = originalDB }()

	auth := NewAuthService(&repository.UserRepository{}, "test-secret")
	if err := auth.Register("pwd-user-2", "password123", "pwd-user-2@example.com"); err != nil {
		t.Fatalf("register user: %v", err)
	}

	_, refreshToken, err := auth.Login("pwd-user-2", "password123")
	if err != nil {
		t.Fatalf("login user: %v", err)
	}

	user, err := (&repository.UserRepository{}).GetByUsername("pwd-user-2")
	if err != nil {
		t.Fatalf("get user: %v", err)
	}

	if err := auth.ChangePassword(user.ID, "password123", "new-password-123"); err != nil {
		t.Fatalf("change password: %v", err)
	}

	if _, _, err := auth.RefreshToken(refreshToken); err == nil {
		t.Fatal("expected old refresh token to be revoked after password change")
	}

	if _, _, err := auth.Login("pwd-user-2", "new-password-123"); err != nil {
		t.Fatalf("expected login with new password to succeed: %v", err)
	}
}

func TestUpdateAvatarPersistsAvatarURL(t *testing.T) {
	db := setupAuthServiceTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() { repository.DB = originalDB }()

	auth := NewAuthService(&repository.UserRepository{}, "test-secret")
	if err := auth.Register("avatar-user", "password123", "avatar-user@example.com"); err != nil {
		t.Fatalf("register user: %v", err)
	}

	user, err := (&repository.UserRepository{}).GetByUsername("avatar-user")
	if err != nil {
		t.Fatalf("get user: %v", err)
	}

	updated, err := auth.UpdateAvatar(user.ID, "/uploads/avatars/1-avatar.png")
	if err != nil {
		t.Fatalf("update avatar: %v", err)
	}

	if updated.AvatarURL != "/uploads/avatars/1-avatar.png" {
		t.Fatalf("expected avatar url to persist, got %s", updated.AvatarURL)
	}
}
