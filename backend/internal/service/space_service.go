package service

import (
	apperror "life-financial-assistant-backend/internal/error"

	"life-financial-assistant-backend/internal/model"
	"life-financial-assistant-backend/internal/repository"
)

type SpaceService struct {
	SpaceRepo *repository.SpaceRepository
	UserRepo  *repository.UserRepository
}

func (s *SpaceService) CreateSpace(userID uint, name, description string) (*model.Space, error) {
	exists, err := s.SpaceRepo.ExistsByName(name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperror.NewConflictError("space name already exists")
	}

	space := &model.Space{
		Name:        name,
		Description: description,
	}

	err = s.SpaceRepo.Create(space)
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

func (s *SpaceService) AddMember(spaceID uint, email string, roleStr string) error {
	user, err := s.UserRepo.GetByEmail(email)
	if err != nil {
		return apperror.NewNotFoundError("user not found")
	}

	role := model.Role(roleStr)

	_, err = s.SpaceRepo.GetUserRoleInSpace(user.ID, spaceID)
	if err == nil {
		return apperror.NewConflictError("user is already a member of this space")
	}

	return s.SpaceRepo.AddUserToSpace(user.ID, spaceID, role)
}

func (s *SpaceService) RemoveMember(operatorID, spaceID, targetUserID uint) error {
	operatorRole, err := s.SpaceRepo.GetUserRoleInSpace(operatorID, spaceID)
	if err != nil {
		return apperror.NewForbiddenError("permission denied")
	}
	if operatorID == targetUserID {
		return apperror.NewForbiddenError("permission denied")
	}

	targetRole, err := s.SpaceRepo.GetUserRoleInSpace(targetUserID, spaceID)
	if err != nil {
		return apperror.NewNotFoundError("user is not a member of this space")
	}

	if targetRole == model.RoleCreator {
		return apperror.NewForbiddenError("permission denied")
	}

	if operatorRole == model.RoleCreator {
		return s.SpaceRepo.RemoveUserFromSpace(targetUserID, spaceID)
	}

	if operatorRole == model.RoleAdmin && (targetRole == model.RoleMember || targetRole == model.RoleObserver) {
		return s.SpaceRepo.RemoveUserFromSpace(targetUserID, spaceID)
	}

	return apperror.NewForbiddenError("permission denied")
}

func (s *SpaceService) UpdateMemberRole(operatorID, spaceID, targetUserID uint, newRole model.Role) error {
	operatorRole, err := s.SpaceRepo.GetUserRoleInSpace(operatorID, spaceID)
	if err != nil {
		return apperror.NewForbiddenError("permission denied")
	}
	if operatorID == targetUserID {
		return apperror.NewForbiddenError("permission denied")
	}

	targetRole, err := s.SpaceRepo.GetUserRoleInSpace(targetUserID, spaceID)
	if err != nil {
		return apperror.NewNotFoundError("user is not a member of this space")
	}

	if targetRole == model.RoleCreator || newRole == model.RoleCreator {
		return apperror.NewForbiddenError("permission denied")
	}

	if operatorRole == model.RoleCreator {
		if newRole == model.RoleAdmin || newRole == model.RoleMember || newRole == model.RoleObserver {
			return s.SpaceRepo.UpdateUserRoleInSpace(targetUserID, spaceID, newRole)
		}
	}

	if operatorRole == model.RoleAdmin {
		if targetRole == model.RoleMember || targetRole == model.RoleObserver {
			if newRole == model.RoleMember || newRole == model.RoleObserver {
				return s.SpaceRepo.UpdateUserRoleInSpace(targetUserID, spaceID, newRole)
			}
		}
	}

	return apperror.NewForbiddenError("permission denied")
}

func (s *SpaceService) LeaveSpace(spaceID, userID uint) error {
	role, err := s.SpaceRepo.GetUserRoleInSpace(userID, spaceID)
	if err != nil {
		return apperror.NewNotFoundError("user is not a member of this space")
	}

	if role == model.RoleCreator {
		return apperror.NewConflictError("creator cannot leave space, please delete instead")
	}

	return s.SpaceRepo.RemoveUserFromSpace(userID, spaceID)
}

func (s *SpaceService) GetMembers(spaceID uint) ([]model.SpaceMember, error) {
	return s.SpaceRepo.GetSpaceMembers(spaceID)
}
