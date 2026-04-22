package repository

import (
	"life-financial-assistant-backend/internal/model"

	"gorm.io/gorm"
)

type BillRepository struct{}

type BillAccessRepository struct{}

func (r BillAccessRepository) GetByID(id uint) (*model.Bill, error) {
	var bill model.Bill
	err := DB.First(&bill, id).Error
	if err != nil {
		return nil, err
	}
	if bill.Amount >= 0 {
		bill.Type = "income"
	} else {
		bill.Type = "expense"
	}
	return &bill, nil
}

func (r *BillRepository) Create(bill *model.Bill) error {
	return DB.Create(bill).Error
}

func applyBillDateRangeFilter(query *gorm.DB, startDate, endDate string) *gorm.DB {
	if startDate != "" {
		query = query.Where("DATE(date) >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("DATE(date) <= ?", endDate)
	}
	return query
}

func (r *BillRepository) GetPersonalBills(userID uint, billType, category, startDate, endDate string) ([]model.Bill, error) {
	var bills []model.Bill
	query := DB.Preload("User").Where("user_id = ? AND is_personal = ?", userID, true)

	if billType != "" {
		if billType == "income" {
			query = query.Where("amount >= ?", 0)
		} else if billType == "expense" {
			query = query.Where("amount < ?", 0)
		}
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	query = applyBillDateRangeFilter(query, startDate, endDate)

	err := query.Find(&bills).Error
	if err != nil {
		return nil, err
	}
	for i := range bills {
		bills[i].CreatorName = bills[i].User.Username
		if bills[i].Amount >= 0 {
			bills[i].Type = "income"
		} else {
			bills[i].Type = "expense"
		}
	}
	return bills, nil
}

func (r *BillRepository) GetSpaceBills(spaceID uint, billType, category, startDate, endDate string) ([]model.Bill, error) {
	var bills []model.Bill
	query := DB.Preload("User").Where("space_id = ? AND is_personal = ?", spaceID, false)

	if billType != "" {
		if billType == "income" {
			query = query.Where("amount >= ?", 0)
		} else if billType == "expense" {
			query = query.Where("amount < ?", 0)
		}
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	query = applyBillDateRangeFilter(query, startDate, endDate)

	err := query.Find(&bills).Error
	if err != nil {
		return nil, err
	}
	for i := range bills {
		bills[i].CreatorName = bills[i].User.Username
		if bills[i].Amount >= 0 {
			bills[i].Type = "income"
		} else {
			bills[i].Type = "expense"
		}
	}
	return bills, nil
}

func (r *BillRepository) GetSharedBills(userID uint, billType, category, startDate, endDate string) ([]model.Bill, error) {
	var bills []model.Bill
	query := DB.Preload("User").
		Joins("JOIN space_user_links on space_user_links.space_id = bills.space_id").
		Where("space_user_links.user_id = ? AND bills.is_personal = ?", userID, false)

	if billType != "" {
		if billType == "income" {
			query = query.Where("bills.amount >= ?", 0)
		} else if billType == "expense" {
			query = query.Where("bills.amount < ?", 0)
		}
	}
	if category != "" {
		query = query.Where("bills.category = ?", category)
	}
	query = applyBillDateRangeFilter(query, startDate, endDate)

	err := query.Find(&bills).Error
	if err != nil {
		return nil, err
	}
	for i := range bills {
		bills[i].CreatorName = bills[i].User.Username
		if bills[i].Amount >= 0 {
			bills[i].Type = "income"
		} else {
			bills[i].Type = "expense"
		}
	}
	return bills, nil
}

func (r *BillRepository) GetByID(id uint) (*model.Bill, error) {
	return BillAccessRepository{}.GetByID(id)
}

func (r *BillRepository) Update(id uint, amount float64, category, billType, remarks string) (*model.Bill, error) {
	var bill model.Bill
	err := DB.First(&bill, id).Error
	if err != nil {
		return nil, err
	}

	bill.Amount = amount
	bill.Category = category
	bill.Type = billType
	bill.Remarks = remarks

	err = DB.Save(&bill).Error
	if err != nil {
		return nil, err
	}

	return &bill, nil
}

func (r *BillRepository) Delete(id uint) error {
	return DB.Delete(&model.Bill{}, id).Error
}