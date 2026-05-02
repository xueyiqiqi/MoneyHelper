package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"life-financial-assistant-backend/internal/model"
	"life-financial-assistant-backend/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo  *repository.UserRepository
	SecretKey []byte
}

func NewAuthService(userRepo *repository.UserRepository, secretKey string) *AuthService {
	return &AuthService{
		UserRepo:  userRepo,
		SecretKey: []byte(secretKey),
	}
}

func (s *AuthService) Register(username, password, email string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Username:       username,
		HashedPassword: string(hashedPassword),
		Email:          email,
	}

	return s.UserRepo.Create(user)
}

func (s *AuthService) Login(username, password string) (string, string, error) {
	user, err := s.UserRepo.GetByUsername(username)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return "", "", errors.New("password mismatch")
	}

	return s.generateTokens(user)
}

func hashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *AuthService) generateTokens(user *model.User) (string, string, error) {
	accessClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
		"type":     "access",
	})
	accessToken, err := accessClaim.SignedString(s.SecretKey)
	if err != nil {
		return "", "", err
	}

	refreshClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
		"type":     "refresh",
	})
	refreshToken, err := refreshClaim.SignedString(s.SecretKey)
	if err != nil {
		return "", "", err
	}

	err = s.UserRepo.UpdateRefreshToken(user.ID, hashRefreshToken(refreshToken))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (string, string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.SecretKey, nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return "", "", errors.New("invalid token type")
	}

	userID := uint(claims["user_id"].(float64))

	user, err := s.UserRepo.GetByID(userID)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	if user.HashedRefreshToken == "" || user.HashedRefreshToken != hashRefreshToken(refreshToken) {
		return "", "", errors.New("refresh token revoked or invalid")
	}

	return s.generateTokens(user)
}

func (s *AuthService) GetProfile(userID uint) (*model.User, error) {
	return s.UserRepo.GetByID(userID)
}

func (s *AuthService) UpdateProfile(userID uint, username, email string) (*model.User, error) {
	exists, err := s.UserRepo.ExistsByUsernameExcludingID(username, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	exists, err = s.UserRepo.ExistsByEmailExcludingID(email, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	if err := s.UserRepo.UpdateProfile(userID, username, email); err != nil {
		return nil, err
	}

	return s.UserRepo.GetByID(userID)
}

func (s *AuthService) UpdateAvatar(userID uint, avatarURL string) (*model.User, error) {
	if err := s.UserRepo.UpdateAvatarURL(userID, avatarURL); err != nil {
		return nil, err
	}

	return s.UserRepo.GetByID(userID)
}

func (s *AuthService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	if currentPassword == newPassword {
		return errors.New("new password cannot be the same as current password")
	}

	user, err := s.UserRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.UserRepo.UpdatePassword(userID, string(hashedPassword)); err != nil {
		return err
	}

	return s.UserRepo.ClearRefreshToken(userID)
}

func (s *AuthService) Logout(userID uint) error {
	return s.UserRepo.ClearRefreshToken(userID)
}
