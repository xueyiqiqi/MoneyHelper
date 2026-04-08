package model

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Username       string         `gorm:"uniqueIndex;not null" json:"username"`
	Email          string         `json:"email"`
	HashedPassword string         `gorm:"not null" json:"-"`
	Spaces         []Space        `gorm:"many2many:space_user_links;" json:"spaces,omitempty"`
}

type Space struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Users       []User         `gorm:"many2many:space_user_links;" json:"users,omitempty"`
	Bills       []Bill         `json:"bills,omitempty"`
}

type SpaceUserLink struct {
	UserID  uint `gorm:"primaryKey"`
	SpaceID uint `gorm:"primaryKey"`
	Role    Role `gorm:"type:varchar(20);default:'member'"`
}
