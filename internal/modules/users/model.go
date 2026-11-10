package users

import "time"

type User struct {
	ID           uint64 `gorm:"primaryKey"`
	Name         string `gorm:"size:120;not null"`
	Email        string `gorm:"size:190;uniqueIndex;not null"`
	PasswordHash string `gorm:"column:password_hash;size:255;not null"`
	IsActive     bool   `gorm:"not null;default:true"`
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (User) TableName() string { return "users" }
