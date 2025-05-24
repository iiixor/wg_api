package models

import (
	"time"
)

type Server struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	Name       string     `json:"name" gorm:"unique;not null"`
	PrivateKey string     `json:"private_key" gorm:"not null"`
	PublicKey  string     `json:"public_key" gorm:"not null"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
