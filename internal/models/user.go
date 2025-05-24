package models

import (
	"time"
)

type User struct {
	ID             uint            `json:"id" gorm:"primaryKey"`
	Username       string          `json:"username" gorm:"unique;not null"`
	ChatID         string          `json:"chat_id" gorm:"unique"`
	Configurations []Configuration `json:"configurations,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      *time.Time      `json:"deleted_at,omitempty" gorm:"index"`
}
