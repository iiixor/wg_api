package repository

import (
	"context"
	"time"
	"wg_api/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
}

type ConfigurationRepository interface {
	Create(ctx context.Context, config *models.Configuration) error
	GetByID(ctx context.Context, id uint) (*models.Configuration, error)
	GetAll(ctx context.Context) ([]models.Configuration, error)
	GetByUserID(ctx context.Context, userID uint) ([]models.Configuration, error)
	GetExpiredConfigs(ctx context.Context) ([]models.Configuration, error)
	GetByStatus(ctx context.Context, status models.ConfigStatus) ([]models.Configuration, error)
	Update(ctx context.Context, config *models.Configuration) error
	UpdateStatus(ctx context.Context, id uint, status models.ConfigStatus) error
	UpdateLatestHandshake(ctx context.Context, id uint, handshakeTime time.Time) error
	Delete(ctx context.Context, id uint) error
}

type ServerRepository interface {
	Create(ctx context.Context, server *models.Server) error
	GetByID(ctx context.Context, id uint) (*models.Server, error)
	GetAll(ctx context.Context) ([]models.Server, error)
	Update(ctx context.Context, server *models.Server) error
	Delete(ctx context.Context, id uint) error
}
