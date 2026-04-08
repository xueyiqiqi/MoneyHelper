package model

import (
	"time"

	"gorm.io/gorm"
)

type Bill struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Amount      float64        `gorm:"not null" json:"amount"`
	Category    string         `gorm:"not null" json:"category"`
	Remarks     string         `json:"remarks"`
	Date        time.Time      `gorm:"index" json:"date"`
	IsPersonal  bool           `gorm:"default:true" json:"is_personal"`
	UserID      uint           `gorm:"index" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"-"` // 后端关联
	CreatorName string         `gorm:"-" json:"creator_name"`      // 仅用于 JSON 返回
	SpaceID     *uint          `gorm:"index" json:"space_id,omitempty"`
}

type AnalysisReport struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Content   string         `gorm:"type:text" json:"content"`
	UserID    uint           `gorm:"index" json:"user_id"`
	SpaceID   *uint          `gorm:"index" json:"space_id,omitempty"`
}
