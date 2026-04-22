package service

import (
	"testing"

	apperror "life-financial-assistant-backend/internal/error"
	"life-financial-assistant-backend/internal/model"
	"life-financial-assistant-backend/internal/repository"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupSpaceServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Space{}, &model.SpaceUserLink{}); err != nil {
		t.Fatalf("migrate schema: %v", err)
	}

	return db
}

func TestAddMemberUsesEmailLookup(t *testing.T) {
	db := setupSpaceServiceTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() {
		repository.DB = originalDB
	}()

	space := model.Space{Name: "email-space"}
	user := model.User{Username: "user-by-email", Email: "invitee@example.com", HashedPassword: "x"}
	if err := db.Create(&space).Error; err != nil {
		t.Fatalf("create space: %v", err)
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	svc := &SpaceService{SpaceRepo: &repository.SpaceRepository{}, UserRepo: &repository.UserRepository{}}

	if err := svc.AddMember(space.ID, user.Email, string(model.RoleMember)); err != nil {
		t.Fatalf("expected add member by email to succeed, got %v", err)
	}

	role, err := (&repository.SpaceRepository{}).GetUserRoleInSpace(user.ID, space.ID)
	if err != nil {
		t.Fatalf("expected user to be added to space, got %v", err)
	}
	if role != model.RoleMember {
		t.Fatalf("expected role member, got %s", role)
	}
}

func TestAddMemberRejectsUnknownEmail(t *testing.T) {
	db := setupSpaceServiceTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() {
		repository.DB = originalDB
	}()

	space := model.Space{Name: "missing-email-space"}
	if err := db.Create(&space).Error; err != nil {
		t.Fatalf("create space: %v", err)
	}

	svc := &SpaceService{SpaceRepo: &repository.SpaceRepository{}, UserRepo: &repository.UserRepository{}}

	err := svc.AddMember(space.ID, "missing@example.com", string(model.RoleMember))
	if err == nil {
		t.Fatal("expected unknown email to return error")
	}

	appErr, ok := err.(*apperror.AppError)
	if !ok || appErr.Code != 404 {
		t.Fatalf("expected 404 app error, got %#v", err)
	}
}


func TestRemoveMemberAllowsCreatorToRemoveAdmin(t *testing.T) {
	db := setupSpaceServiceTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() { repository.DB = originalDB }()

	creator := model.User{Username: "creator-remove", Email: "creator-remove@example.com", HashedPassword: "x"}
	admin := model.User{Username: "admin-remove", Email: "admin-remove@example.com", HashedPassword: "x"}
	space := model.Space{Name: "remove-space"}
	_ = db.Create(&creator).Error
	_ = db.Create(&admin).Error
	_ = db.Create(&space).Error
	_ = db.Create(&model.SpaceUserLink{UserID: creator.ID, SpaceID: space.ID, Role: model.RoleCreator}).Error
	_ = db.Create(&model.SpaceUserLink{UserID: admin.ID, SpaceID: space.ID, Role: model.RoleAdmin}).Error

	svc := &SpaceService{SpaceRepo: &repository.SpaceRepository{}, UserRepo: &repository.UserRepository{}}
	if err := svc.RemoveMember(creator.ID, space.ID, admin.ID); err != nil {
		t.Fatalf("expected creator remove admin to succeed, got %v", err)
	}
}

func TestRemoveMemberRejectsAdminRemovingAdmin(t *testing.T) {
	db := setupSpaceServiceTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() { repository.DB = originalDB }()

	admin1 := model.User{Username: "admin-1", Email: "admin-1@example.com", HashedPassword: "x"}
	admin2 := model.User{Username: "admin-2", Email: "admin-2@example.com", HashedPassword: "x"}
	space := model.Space{Name: "remove-space-2"}
	_ = db.Create(&admin1).Error
	_ = db.Create(&admin2).Error
	_ = db.Create(&space).Error
	_ = db.Create(&model.SpaceUserLink{UserID: admin1.ID, SpaceID: space.ID, Role: model.RoleAdmin}).Error
	_ = db.Create(&model.SpaceUserLink{UserID: admin2.ID, SpaceID: space.ID, Role: model.RoleAdmin}).Error

	svc := &SpaceService{SpaceRepo: &repository.SpaceRepository{}, UserRepo: &repository.UserRepository{}}
	err := svc.RemoveMember(admin1.ID, space.ID, admin2.ID)
	if err == nil {
		t.Fatal("expected admin removing admin to fail")
	}
	appErr, ok := err.(*apperror.AppError)
	if !ok || appErr.Code != 403 {
		t.Fatalf("expected 403 app error, got %#v", err)
	}
}

func TestUpdateMemberRoleAllowsCreatorToPromoteAdmin(t *testing.T) {
	db := setupSpaceServiceTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() { repository.DB = originalDB }()

	creator := model.User{Username: "creator-promote", Email: "creator-promote@example.com", HashedPassword: "x"}
	member := model.User{Username: "member-promote", Email: "member-promote@example.com", HashedPassword: "x"}
	space := model.Space{Name: "role-space"}
	_ = db.Create(&creator).Error
	_ = db.Create(&member).Error
	_ = db.Create(&space).Error
	_ = db.Create(&model.SpaceUserLink{UserID: creator.ID, SpaceID: space.ID, Role: model.RoleCreator}).Error
	_ = db.Create(&model.SpaceUserLink{UserID: member.ID, SpaceID: space.ID, Role: model.RoleMember}).Error

	svc := &SpaceService{SpaceRepo: &repository.SpaceRepository{}, UserRepo: &repository.UserRepository{}}
	if err := svc.UpdateMemberRole(creator.ID, space.ID, member.ID, model.RoleAdmin); err != nil {
		t.Fatalf("expected creator promote admin to succeed, got %v", err)
	}
}

func TestUpdateMemberRoleRejectsAdminPromotingAdmin(t *testing.T) {
	db := setupSpaceServiceTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() { repository.DB = originalDB }()

	admin := model.User{Username: "admin-promoter", Email: "admin-promoter@example.com", HashedPassword: "x"}
	member := model.User{Username: "member-promoted", Email: "member-promoted@example.com", HashedPassword: "x"}
	space := model.Space{Name: "role-space-2"}
	_ = db.Create(&admin).Error
	_ = db.Create(&member).Error
	_ = db.Create(&space).Error
	_ = db.Create(&model.SpaceUserLink{UserID: admin.ID, SpaceID: space.ID, Role: model.RoleAdmin}).Error
	_ = db.Create(&model.SpaceUserLink{UserID: member.ID, SpaceID: space.ID, Role: model.RoleMember}).Error

	svc := &SpaceService{SpaceRepo: &repository.SpaceRepository{}, UserRepo: &repository.UserRepository{}}
	err := svc.UpdateMemberRole(admin.ID, space.ID, member.ID, model.RoleAdmin)
	if err == nil {
		t.Fatal("expected admin promote admin to fail")
	}
	appErr, ok := err.(*apperror.AppError)
	if !ok || appErr.Code != 403 {
		t.Fatalf("expected 403 app error, got %#v", err)
	}
}


func TestGetUserSpacesIncludesCurrentUserRole(t *testing.T) {
	db := setupSpaceServiceTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() { repository.DB = originalDB }()

	user := model.User{Username: "role-user", Email: "role-user@example.com", HashedPassword: "x"}
	space := model.Space{Name: "role-visible-space"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := db.Create(&space).Error; err != nil {
		t.Fatalf("create space: %v", err)
	}
	if err := db.Create(&model.SpaceUserLink{UserID: user.ID, SpaceID: space.ID, Role: model.RoleCreator}).Error; err != nil {
		t.Fatalf("create link: %v", err)
	}

	svc := &SpaceService{SpaceRepo: &repository.SpaceRepository{}, UserRepo: &repository.UserRepository{}}
	spaces, err := svc.GetUserSpaces(user.ID)
	if err != nil {
		t.Fatalf("get user spaces: %v", err)
	}
	if len(spaces) != 1 {
		t.Fatalf("expected 1 space, got %d", len(spaces))
	}
	if spaces[0].Role != model.RoleCreator {
		t.Fatalf("expected role creator, got %q", spaces[0].Role)
	}
}
