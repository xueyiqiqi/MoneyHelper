package service

import (
	"life-financial-assistant-backend/internal/model"
	"life-financial-assistant-backend/internal/repository"
)

type SpaceService struct {
	SpaceRepo *repository.SpaceRepository
}

func (s *SpaceService) CreateSpace(userID uint, name, description string) (*model.Space, error) {
	space := &model.Space{
		Name:        name,
		Description: description,
	}

	err := s.SpaceRepo.Create(space)
	if err != nil {
		return nil, err
	}

	err = s.SpaceRepo.AddUserToSpace(userID, space.ID, model.RoleCreator)
	if err != nil {
		return nil, err
	}

	return space, nil
}

func (s *SpaceService) GetUserSpaces(userID uint) ([]model.Space, error) {
	return s.SpaceRepo.GetUserSpaces(userID)
}

func (s *SpaceService) CheckPermission(userID, spaceID uint, requiredRoles []model.Role) (bool, error) {
	role, err := s.SpaceRepo.GetUserRoleInSpace(userID, spaceID)
	if err != nil {
		return false, err
	}

	for _, r := range requiredRoles {
		if r == role {
			return true, nil
		}
	}

	return false, nil
}
