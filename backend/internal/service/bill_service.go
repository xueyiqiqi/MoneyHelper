package service

import (
	apperror "life-financial-assistant-backend/internal/error"
	"life-financial-assistant-backend/internal/model"
	"life-financial-assistant-backend/internal/repository"
	"life-financial-assistant-backend/pkg/ai"
	"time"
)

type BillService struct {
	BillRepo   *repository.BillRepository
	SpaceRepo  *repository.SpaceRepository
	ReportRepo *repository.ReportRepository
	AIAgent    *ai.AIAgent
}

func (s *BillService) CreateBill(userID uint, amount float64, category, billType, remarks string, isPersonal bool, spaceID *uint) (*model.Bill, error) {
	// 如果是支出，金额存储为负数
	if billType == "expense" {
		amount = -amount
	}

	bill := &model.Bill{
		Amount:     amount,
		Category:   category,
		Type:       billType,
		Remarks:    remarks,
		Date:       time.Now(),
		IsPersonal: isPersonal,
		UserID:     userID,
		SpaceID:    spaceID,
	}

	return bill, s.BillRepo.Create(bill)
}

func (s *BillService) GetUserBills(userID uint, isPersonal bool, spaceID *uint, billType, category, startDate, endDate string) ([]model.Bill, error) {
	if isPersonal {
		return s.BillRepo.GetPersonalBills(userID, billType, category, startDate, endDate)
	}
	if spaceID != nil {
		return s.BillRepo.GetSpaceBills(*spaceID, billType, category, startDate, endDate)
	}
	return s.BillRepo.GetSharedBills(userID, billType, category, startDate, endDate)
}

func (s *BillService) GetBillByID(userID, id uint) (*model.Bill, error) {
	bill, err := s.BillRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if bill.IsPersonal {
		if bill.UserID != userID {
			return nil, apperror.NewForbiddenError("permission denied")
		}
		return bill, nil
	}

	if bill.SpaceID == nil {
		return nil, apperror.NewForbiddenError("permission denied")
	}

	_, err = s.SpaceRepo.GetUserRoleInSpace(userID, *bill.SpaceID)
	if err != nil {
		return nil, apperror.NewForbiddenError("permission denied")
	}

	return bill, nil
}

func (s *BillService) UpdateBill(userID uint, id uint, amount float64, category, billType, remarks string, isPersonal bool, spaceID *uint) (*model.Bill, error) {
	if _, err := s.GetBillByID(userID, id); err != nil {
		return nil, err
	}

	if billType == "expense" {
		amount = -amount
	}

	return s.BillRepo.Update(id, amount, category, billType, remarks)
}

func (s *BillService) GenerateAnalysis(userID uint, isPersonal bool, spaceID *uint) (*model.AnalysisReport, error) {
	bills, err := s.GetUserBills(userID, isPersonal, spaceID, "", "", "", "")
	if err != nil {
		return nil, err
	}

	context := "personal"
	if spaceID != nil {
		context = "space"
	} else if !isPersonal {
		context = "shared"
	}

	content, err := s.AIAgent.AnalyzeBills(bills, context)
	if err != nil {
		return nil, err
	}

	report := &model.AnalysisReport{
		Content: content,
		UserID:  userID,
		SpaceID: spaceID,
	}

	err = s.ReportRepo.Create(report)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (s *BillService) GetReports(userID uint, isPersonal bool, spaceID *uint) ([]model.AnalysisReport, error) {
	if isPersonal {
		return s.ReportRepo.GetPersonalReports(userID)
	}
	if spaceID != nil {
		return s.ReportRepo.GetSpaceReports(*spaceID)
	}
	return s.ReportRepo.GetSharedReports(userID)
}
