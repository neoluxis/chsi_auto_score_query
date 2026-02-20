package model

import (
"time"

"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey"`
	Name         string         `gorm:"index"`
	IDCard       string         `gorm:"index"`
	ExamID       string         `gorm:"index"`
	Email        string
	SchoolCode   string
	InfoHash     string    `gorm:"uniqueIndex"`
	Score        string    `gorm:"type:text"`
	Notice       string    `gorm:"type:text"`
	LastQueryAt  time.Time `gorm:"index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
