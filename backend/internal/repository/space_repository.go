package repository

import (
	"life-financial-assistant-backend/internal/model"
)

type SpaceRepository struct{}

func (r *SpaceRepository) Create(space *model.Space) error {
	return DB.Create(space).Error
}

func (r *SpaceRepository) AddUserToSpace(userID, spaceID uint, role model.Role) error {
	link := model.SpaceUserLink{
		UserID:  userID,
		SpaceID: spaceID,
		Role:    role,
	}
	return DB.Create(&link).Error
}

func (r *SpaceRepository) GetUserRoleInSpace(userID, spaceID uint) (model.Role, error) {
	var link model.SpaceUserLink
	err := DB.Where("user_id = ? AND space_id = ?", userID, spaceID).First(&link).Error
	if err != nil {
		return "", err
	}
	return link.Role, nil
}

func (r *SpaceRepository) GetUserSpaces(userID uint) ([]model.Space, error) {
	var user model.User
	err := DB.Preload("Spaces").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return user.Spaces, nil
}

func (r *SpaceRepository) GetByID(id uint) (*model.Space, error) {
	var space model.Space
	err := DB.First(&space, id).Error
	if err != nil {
		return nil, err
	}
	return &space, nil
}
