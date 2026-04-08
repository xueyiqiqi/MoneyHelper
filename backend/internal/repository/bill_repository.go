package repository

import (
	"life-financial-assistant-backend/internal/model"
)

type BillRepository struct{}

func (r *BillRepository) Create(bill *model.Bill) error {
	return DB.Create(bill).Error
}

func (r *BillRepository) GetPersonalBills(userID uint) ([]model.Bill, error) {
	var bills []model.Bill
	err := DB.Preload("User").Where("user_id = ? AND is_personal = ?", userID, true).Find(&bills).Error
	if err != nil {
		return nil, err
	}
	for i := range bills {
		bills[i].CreatorName = bills[i].User.Username
	}
	return bills, nil
}

func (r *BillRepository) GetSpaceBills(spaceID uint) ([]model.Bill, error) {
	var bills []model.Bill
	err := DB.Preload("User").Where("space_id = ? AND is_personal = ?", spaceID, false).Find(&bills).Error
	if err != nil {
		return nil, err
	}
	for i := range bills {
		bills[i].CreatorName = bills[i].User.Username
	}
	return bills, nil
}

func (r *BillRepository) GetSharedBills(userID uint) ([]model.Bill, error) {
	var bills []model.Bill
	// 查找用户所属的所有空间的账单
	err := DB.Preload("User").
		Joins("JOIN space_user_links on space_user_links.space_id = bills.space_id").
		Where("space_user_links.user_id = ? AND bills.is_personal = ?", userID, false).
		Find(&bills).Error
	if err != nil {
		return nil, err
	}
	for i := range bills {
		bills[i].CreatorName = bills[i].User.Username
	}
	return bills, nil
}

func (r *BillRepository) GetByID(id uint) (*model.Bill, error) {
	var bill model.Bill
	err := DB.First(&bill, id).Error
	if err != nil {
		return nil, err
	}
	return &bill, nil
}

func (r *BillRepository) Delete(id uint) error {
	return DB.Delete(&model.Bill{}, id).Error
}
