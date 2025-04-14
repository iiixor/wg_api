package models

import (
	"fmt"
	"time"
)

type ConfigStatus string

const (
	StatusPaid     ConfigStatus = "paid"
	StatusExpired  ConfigStatus = "expired"
	StatusDeletion ConfigStatus = "deletion"
)

type Configuration struct {
	ID              uint         `json:"id" gorm:"primaryKey"`
	Name            string       `json:"name" gorm:"not null"`
	Status          ConfigStatus `json:"status" gorm:"type:varchar(20);default:'new'"`
	ExpirationTime  time.Time    `json:"expiration_time"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	DeletedAt       *time.Time   `json:"deleted_at,omitempty" gorm:"index"`
	InterfaceID     uint         `json:"interface_id" gorm:"not null"` // Сервер
	ServerID        uint         `json:"-" gorm:"column:interface_id"` // Алиас для InterfaceID
	PrivateKey      string       `json:"private_key" gorm:"not null"`
	PublicKey       string       `json:"public_key" gorm:"not null"`
	AllowedIP       string       `json:"allowed_ip" gorm:"not null"`
	LatestHandshake *time.Time   `json:"latest_handshake"`
	UserID          uint         `json:"user_id" gorm:"not null"`
	Server          *Server      `json:"server,omitempty" gorm:"foreignKey:InterfaceID"`
	User            *User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (c *Configuration) ToWireGuardConfig() string {
	return fmt.Sprintf(
		"[Peer]\n"+
			"# Name = %s\n"+
			"PublicKey = %s\n"+
			"AllowedIPs = %s\n\n",
		c.Name,
		c.PublicKey,
		c.AllowedIP)
}
