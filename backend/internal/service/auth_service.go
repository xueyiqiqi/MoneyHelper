package service

import (
	"errors"
	"life-financial-assistant-backend/internal/model"
	"life-financial-assistant-backend/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	SecretKey = []byte("your-secret-key")
)

type AuthService struct {
	UserRepo *repository.UserRepository
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

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.UserRepo.GetByUsername(username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(SecretKey)
}
