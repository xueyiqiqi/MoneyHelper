package service

import (
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

func (s *BillService) CreateBill(userID uint, amount float64, category, remarks string, isPersonal bool, spaceID *uint) (*model.Bill, error) {
	bill := &model.Bill{
		Amount:     amount,
		Category:   category,
		Remarks:    remarks,
		Date:       time.Now(),
		IsPersonal: isPersonal,
		UserID:     userID,
		SpaceID:    spaceID,
	}

	return bill, s.BillRepo.Create(bill)
}

func (s *BillService) GetUserBills(userID uint, isPersonal bool, spaceID *uint) ([]model.Bill, error) {
	if isPersonal {
		return s.BillRepo.GetPersonalBills(userID)
	}
	if spaceID != nil {
		return s.BillRepo.GetSpaceBills(*spaceID)
	}
	return s.BillRepo.GetSharedBills(userID)
}

func (s *BillService) GenerateAnalysis(userID uint, isPersonal bool, spaceID *uint) (*model.AnalysisReport, error) {
	bills, err := s.GetUserBills(userID, isPersonal, spaceID)
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
