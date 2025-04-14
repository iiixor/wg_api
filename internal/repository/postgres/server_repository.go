package postgres

import (
	"context"
	"errors"
	"wg_api/internal/models"

	"gorm.io/gorm"
)

type ServerRepository struct {
	db *gorm.DB
}

func NewServerRepository(db *gorm.DB) *ServerRepository {
	return &ServerRepository{db: db}
}

func (r *ServerRepository) Create(ctx context.Context, server *models.Server) error {
	return r.db.WithContext(ctx).Create(server).Error
}

func (r *ServerRepository) GetByID(ctx context.Context, id uint) (*models.Server, error) {
	var server models.Server
	err := r.db.WithContext(ctx).First(&server, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("server not found")
		}
		return nil, err
	}
	return &server, nil
}

func (r *ServerRepository) GetAll(ctx context.Context) ([]models.Server, error) {
	var servers []models.Server
	err := r.db.WithContext(ctx).Find(&servers).Error
	return servers, err
}

func (r *ServerRepository) Update(ctx context.Context, server *models.Server) error {
	return r.db.WithContext(ctx).Save(server).Error
}

func (r *ServerRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Server{}, id).Error
}
