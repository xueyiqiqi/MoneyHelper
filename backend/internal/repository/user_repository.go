package repository

import (
	"life-financial-assistant-backend/internal/model"
)

type UserRepository struct{}

func (r *UserRepository) Create(user *model.User) error {
	return DB.Create(user).Error
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateRefreshToken(userID uint, hashedToken string) error {
	return DB.Model(&model.User{}).Where("id = ?", userID).Update("hashed_refresh_token", hashedToken).Error
}

func (r *UserRepository) ClearRefreshToken(userID uint) error {
	return DB.Model(&model.User{}).Where("id = ?", userID).Update("hashed_refresh_token", "").Error
}
