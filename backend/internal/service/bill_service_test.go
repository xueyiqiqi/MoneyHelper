package service

import (
	"testing"
	"time"

	apperror "life-financial-assistant-backend/internal/error"
	"life-financial-assistant-backend/internal/model"
	"life-financial-assistant-backend/internal/repository"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupBillServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Space{}, &model.SpaceUserLink{}, &model.Bill{}); err != nil {
		t.Fatalf("migrate schema: %v", err)
	}

	return db
}

func TestGetBillByIDRejectsUnauthorizedPersonalBill(t *testing.T) {
	db := setupBillServiceTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() {
		repository.DB = originalDB
	}()

	owner := model.User{Username: "owner", Email: "owner@example.com", HashedPassword: "x"}
	viewer := model.User{Username: "viewer", Email: "viewer@example.com", HashedPassword: "x"}
	if err := db.Create(&owner).Error; err != nil {
		t.Fatalf("create owner: %v", err)
	}
	if err := db.Create(&viewer).Error; err != nil {
		t.Fatalf("create viewer: %v", err)
	}

	bill := model.Bill{
		Amount:     10,
		Category:   "food",
		Type:       "income",
		Remarks:    "test",
		Date:       time.Now(),
		IsPersonal: true,
		UserID:     owner.ID,
	}
	if err := db.Create(&bill).Error; err != nil {
		t.Fatalf("create bill: %v", err)
	}

	svc := &BillService{BillRepo: &repository.BillRepository{}, SpaceRepo: &repository.SpaceRepository{}}

	_, err := svc.GetBillByID(viewer.ID, bill.ID)
	if err == nil {
		t.Fatal("expected unauthorized access to return error")
	}

	appErr, ok := err.(*apperror.AppError)
	if !ok || appErr.Code != 403 {
		t.Fatalf("expected 403 app error, got %#v", err)
	}
}

func TestGetBillByIDAllowsAuthorizedSpaceMember(t *testing.T) {
	db := setupBillServiceTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() {
		repository.DB = originalDB
	}()

	owner := model.User{Username: "space-owner", Email: "space-owner@example.com", HashedPassword: "x"}
	member := model.User{Username: "space-member", Email: "space-member@example.com", HashedPassword: "x"}
	space := model.Space{Name: "shared-space"}
	if err := db.Create(&owner).Error; err != nil {
		t.Fatalf("create owner: %v", err)
	}
	if err := db.Create(&member).Error; err != nil {
		t.Fatalf("create member: %v", err)
	}
	if err := db.Create(&space).Error; err != nil {
		t.Fatalf("create space: %v", err)
	}
	if err := db.Create(&model.SpaceUserLink{UserID: owner.ID, SpaceID: space.ID, Role: model.RoleCreator}).Error; err != nil {
		t.Fatalf("create owner link: %v", err)
	}
	if err := db.Create(&model.SpaceUserLink{UserID: member.ID, SpaceID: space.ID, Role: model.RoleMember}).Error; err != nil {
		t.Fatalf("create member link: %v", err)
	}
	if _, err := (&repository.SpaceRepository{}).GetUserRoleInSpace(member.ID, space.ID); err != nil {
		t.Fatalf("expected member link lookup to succeed, got %v", err)
	}

	bill := model.Bill{
		Amount:     10,
		Category:   "food",
		Type:       "income",
		Remarks:    "shared",
		Date:       time.Now(),
		IsPersonal: false,
		UserID:     owner.ID,
		SpaceID:    &space.ID,
	}
	if err := db.Create(&bill).Error; err != nil {
		t.Fatalf("create bill: %v", err)
	}
	stored, err := (&repository.BillRepository{}).GetByID(bill.ID)
	if err != nil {
		t.Fatalf("expected stored bill lookup to succeed, got %v", err)
	}
	if stored.SpaceID == nil || *stored.SpaceID != space.ID || stored.IsPersonal {
		t.Fatalf("unexpected stored bill state: %#v", stored)
	}

	svc := &BillService{BillRepo: &repository.BillRepository{}, SpaceRepo: &repository.SpaceRepository{}}

	got, err := svc.GetBillByID(member.ID, bill.ID)
	if err != nil {
		t.Fatalf("expected shared-space member access to succeed, got %v", err)
	}
	if got == nil || got.ID != bill.ID {
		t.Fatalf("expected to get shared bill, got %#v", got)
	}
}

func TestGetBillByIDRejectsNonMemberForSharedSpaceBill(t *testing.T) {
	db := setupBillServiceTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() {
		repository.DB = originalDB
	}()

	owner := model.User{Username: "space-owner-2", Email: "space-owner-2@example.com", HashedPassword: "x"}
	viewer := model.User{Username: "space-outsider", Email: "space-outsider@example.com", HashedPassword: "x"}
	space := model.Space{Name: "shared-space-2"}
	if err := db.Create(&owner).Error; err != nil {
		t.Fatalf("create owner: %v", err)
	}
	if err := db.Create(&viewer).Error; err != nil {
		t.Fatalf("create viewer: %v", err)
	}
	if err := db.Create(&space).Error; err != nil {
		t.Fatalf("create space: %v", err)
	}
	if err := db.Create(&model.SpaceUserLink{UserID: owner.ID, SpaceID: space.ID, Role: model.RoleCreator}).Error; err != nil {
		t.Fatalf("create owner link: %v", err)
	}

	bill := model.Bill{
		Amount:     10,
		Category:   "food",
		Type:       "income",
		Remarks:    "shared",
		Date:       time.Now(),
		IsPersonal: false,
		UserID:     owner.ID,
		SpaceID:    &space.ID,
	}
	if err := db.Create(&bill).Error; err != nil {
		t.Fatalf("create bill: %v", err)
	}

	svc := &BillService{BillRepo: &repository.BillRepository{}, SpaceRepo: &repository.SpaceRepository{}}

	_, err := svc.GetBillByID(viewer.ID, bill.ID)
	if err == nil {
		t.Fatal("expected non-member access to shared bill to return error")
	}

	appErr, ok := err.(*apperror.AppError)
	if !ok || appErr.Code != 403 {
		t.Fatalf("expected 403 app error, got %#v", err)
	}
}

func TestGetUserBillsFiltersPersonalBillsByDateRange(t *testing.T) {
	db := setupBillServiceTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() {
		repository.DB = originalDB
	}()

	user := model.User{Username: "date-filter-user", Email: "date-filter-user@example.com", HashedPassword: "x"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	bills := []model.Bill{
		{
			Amount:     10,
			Category:   "food",
			Type:       "income",
			Remarks:    "before range",
			Date:       time.Date(2026, time.April, 9, 8, 0, 0, 0, time.UTC),
			IsPersonal: true,
			UserID:     user.ID,
		},
		{
			Amount:     20,
			Category:   "food",
			Type:       "income",
			Remarks:    "range start",
			Date:       time.Date(2026, time.April, 10, 8, 0, 0, 0, time.UTC),
			IsPersonal: true,
			UserID:     user.ID,
		},
		{
			Amount:     30,
			Category:   "food",
			Type:       "income",
			Remarks:    "in range",
			Date:       time.Date(2026, time.April, 11, 8, 0, 0, 0, time.UTC),
			IsPersonal: true,
			UserID:     user.ID,
		},
		{
			Amount:     40,
			Category:   "food",
			Type:       "income",
			Remarks:    "range end",
			Date:       time.Date(2026, time.April, 12, 8, 0, 0, 0, time.UTC),
			IsPersonal: true,
			UserID:     user.ID,
		},
		{
			Amount:     50,
			Category:   "food",
			Type:       "income",
			Remarks:    "after range",
			Date:       time.Date(2026, time.April, 13, 8, 0, 0, 0, time.UTC),
			IsPersonal: true,
			UserID:     user.ID,
		},
		{
			Amount:     60,
			Category:   "food",
			Type:       "income",
			Remarks:    "shared bill",
			Date:       time.Date(2026, time.April, 11, 8, 0, 0, 0, time.UTC),
			IsPersonal: false,
			UserID:     user.ID,
		},
		{
			Amount:     70,
			Category:   "transport",
			Type:       "income",
			Remarks:    "different category",
			Date:       time.Date(2026, time.April, 11, 8, 0, 0, 0, time.UTC),
			IsPersonal: true,
			UserID:     user.ID,
		},
	}
	if err := db.Create(&bills).Error; err != nil {
		t.Fatalf("create bills: %v", err)
	}

	svc := &BillService{BillRepo: &repository.BillRepository{}, SpaceRepo: &repository.SpaceRepository{}}

	got, err := svc.GetUserBills(user.ID, true, nil, "income", "food", "2026-04-10", "2026-04-12")
	if err != nil {
		t.Fatalf("GetUserBills returned error: %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 personal bills in inclusive range, got %d: %#v", len(got), got)
	}

	for _, bill := range got {
		if !bill.IsPersonal {
			t.Fatalf("expected only personal bills, got %#v", bill)
		}
		if bill.Category != "food" {
			t.Fatalf("expected only food bills, got %#v", bill)
		}
		if bill.Date.Before(time.Date(2026, time.April, 10, 0, 0, 0, 0, time.UTC)) || bill.Date.After(time.Date(2026, time.April, 12, 23, 59, 59, 0, time.UTC)) {
			t.Fatalf("expected bill date within inclusive range, got %#v", bill)
		}
	}
}
