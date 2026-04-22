package repository

import (
	"life-financial-assistant-backend/internal/model"
)

type SpaceRepository struct{}

func (r *SpaceRepository) Create(space *model.Space) error {
	return DB.Create(space).Error
}

func (r *SpaceRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := DB.Model(&model.Space{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
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

	for i := range user.Spaces {
		var link model.SpaceUserLink
		if err := DB.Where("user_id = ? AND space_id = ?", userID, user.Spaces[i].ID).First(&link).Error; err == nil {
			user.Spaces[i].Role = link.Role
		}

		var count int64
		DB.Model(&model.SpaceUserLink{}).Where("space_id = ?", user.Spaces[i].ID).Count(&count)
		user.Spaces[i].MemberCount = int(count)
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

func (r *SpaceRepository) RemoveUserFromSpace(userID, spaceID uint) error {
	return DB.Where("user_id = ? AND space_id = ?", userID, spaceID).Delete(&model.SpaceUserLink{}).Error
}

func (r *SpaceRepository) UpdateUserRoleInSpace(userID, spaceID uint, role model.Role) error {
	return DB.Model(&model.SpaceUserLink{}).Where("user_id = ? AND space_id = ?", userID, spaceID).Update("role", role).Error
}

func (r *SpaceRepository) GetSpaceMembers(spaceID uint) ([]model.SpaceMember, error) {
	var links []model.SpaceUserLink
	err := DB.Preload("User").Where("space_id = ?", spaceID).Find(&links).Error
	if err != nil {
		return nil, err
	}

	members := make([]model.SpaceMember, len(links))
	for i, link := range links {
		members[i] = model.SpaceMember{
			UserID:   link.UserID,
			Username: link.User.Username,
			Role:     link.Role,
		}
	}
	return members, nil
}
