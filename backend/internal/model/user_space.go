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
	Username       string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email          string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	HashedPassword     string         `gorm:"type:varchar(255);not null" json:"-"`
	HashedRefreshToken string         `gorm:"type:text" json:"-"`
	Spaces             []Space        `gorm:"many2many:space_user_links;" json:"spaces,omitempty"`
}

type Space struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:varchar(255)" json:"description"`
	Users       []User         `gorm:"many2many:space_user_links;" json:"users,omitempty"`
	Bills       []Bill         `json:"bills,omitempty"`
	MemberCount int            `gorm:"-" json:"member_count"`
	Role        Role           `gorm:"-" json:"role"`
}

type SpaceUserLink struct {
	UserID  uint `gorm:"primaryKey"`
	SpaceID uint `gorm:"primaryKey"`
	Role    Role `gorm:"type:varchar(20);default:'member'"`
	User    User `gorm:"foreignKey:UserID"`
}

type SpaceMember struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     Role   `json:"role"`
}
