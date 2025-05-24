package postgres

import (
	"context"
	"errors"
	"time"
	"wg_api/internal/models"

	"gorm.io/gorm"
)

type ConfigurationRepository struct {
	db *gorm.DB
}

func NewConfigurationRepository(db *gorm.DB) *ConfigurationRepository {
	return &ConfigurationRepository{db: db}
}

func (r *ConfigurationRepository) Create(ctx context.Context, config *models.Configuration) error {
	return r.db.WithContext(ctx).Create(config).Error
}

func (r *ConfigurationRepository) GetByID(ctx context.Context, id uint) (*models.Configuration, error) {
	var config models.Configuration
	err := r.db.WithContext(ctx).Preload("Server").Preload("User").First(&config, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("configuration not found")
		}
		return nil, err
	}
	return &config, nil
}

func (r *ConfigurationRepository) GetAll(ctx context.Context) ([]models.Configuration, error) {
	var configs []models.Configuration
	err := r.db.WithContext(ctx).Preload("Server").Preload("User").Find(&configs).Error
	return configs, err
}

func (r *ConfigurationRepository) GetByUserID(ctx context.Context, userID uint) ([]models.Configuration, error) {
	var configs []models.Configuration
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Server").Find(&configs).Error
	return configs, err
}

func (r *ConfigurationRepository) GetExpiredConfigs(ctx context.Context) ([]models.Configuration, error) {
	var configs []models.Configuration
	now := time.Now()
	err := r.db.WithContext(ctx).
		Where("status = ? AND expiration_time < ?", models.StatusPaid, now).
		Preload("Server").
		Find(&configs).Error
	return configs, err
}

func (r *ConfigurationRepository) GetByStatus(ctx context.Context, status models.ConfigStatus) ([]models.Configuration, error) {
	var configs []models.Configuration
	err := r.db.WithContext(ctx).Where("status = ?", status).Preload("Server").Find(&configs).Error
	return configs, err
}

func (r *ConfigurationRepository) Update(ctx context.Context, config *models.Configuration) error {
	return r.db.WithContext(ctx).Save(config).Error
}

func (r *ConfigurationRepository) UpdateStatus(ctx context.Context, id uint, status models.ConfigStatus) error {
	return r.db.WithContext(ctx).Model(&models.Configuration{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *ConfigurationRepository) UpdateLatestHandshake(ctx context.Context, id uint, handshakeTime time.Time) error {
	return r.db.WithContext(ctx).Model(&models.Configuration{}).
		Where("id = ?", id).
		Update("latest_handshake", handshakeTime).Error
}

func (r *ConfigurationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Configuration{}, id).Error
}
