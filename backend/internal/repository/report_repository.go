package repository

import (
	"life-financial-assistant-backend/internal/model"
)

type ReportRepository struct{}

func (r *ReportRepository) Create(report *model.AnalysisReport) error {
	return DB.Create(report).Error
}

func (r *ReportRepository) GetPersonalReports(userID uint) ([]model.AnalysisReport, error) {
	var reports []model.AnalysisReport
	err := DB.Where("user_id = ? AND space_id IS NULL", userID).Find(&reports).Error
	if err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *ReportRepository) GetSpaceReports(spaceID uint) ([]model.AnalysisReport, error) {
	var reports []model.AnalysisReport
	err := DB.Where("space_id = ?", spaceID).Find(&reports).Error
	if err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *ReportRepository) GetSharedReports(userID uint) ([]model.AnalysisReport, error) {
	var reports []model.AnalysisReport
	err := DB.Joins("JOIN space_user_links on space_user_links.space_id = analysis_reports.space_id").
		Where("space_user_links.user_id = ? AND analysis_reports.space_id IS NOT NULL", userID).
		Find(&reports).Error
	if err != nil {
		return nil, err
	}
	return reports, nil
}
