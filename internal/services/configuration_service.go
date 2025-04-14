package services

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
	"wg_api/internal/models"
	"wg_api/internal/repository"
)

type ConfigurationService struct {
	repo repository.ConfigurationRepository
}

func NewConfigurationService(repo repository.ConfigurationRepository) *ConfigurationService {
	return &ConfigurationService{repo: repo}
}

func (s *ConfigurationService) Create(ctx context.Context, config *models.Configuration) error {
	// Генерируем ключевую пару
	privateKey, publicKey, err := generateWireGuardKeyPair()
	if err != nil {
		return fmt.Errorf("failed to generate wireguard key pair: %w", err)
	}

	// Устанавливаем значения по умолчанию
	config.PrivateKey = privateKey
	config.PublicKey = publicKey
	config.CreatedAt = time.Now()
	config.ExpirationTime = time.Now().AddDate(0, 1, 0) // +1 месяц
	config.Status = models.StatusPaid                   // По умолчанию "new"

	return s.repo.Create(ctx, config)
}

func (s *ConfigurationService) GetByID(ctx context.Context, id uint) (*models.Configuration, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ConfigurationService) GetAll(ctx context.Context) ([]models.Configuration, error) {
	return s.repo.GetAll(ctx)
}

func (s *ConfigurationService) GetByUserID(ctx context.Context, userID uint) ([]models.Configuration, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *ConfigurationService) GetAllByStatus(ctx context.Context, status models.ConfigStatus) ([]models.Configuration, error) {
	return s.repo.GetByStatus(ctx, status)
}

func (s *ConfigurationService) Update(ctx context.Context, config *models.Configuration) error {
	return s.repo.Update(ctx, config)
}

func (s *ConfigurationService) UpdateStatus(ctx context.Context, id uint, status models.ConfigStatus) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *ConfigurationService) UpdateLatestHandshake(ctx context.Context, id uint, handshakeTime time.Time) error {
	return s.repo.UpdateLatestHandshake(ctx, id, handshakeTime)
}

func (s *ConfigurationService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *ConfigurationService) UpdateExpiredConfigurations(ctx context.Context) error {
	// Получаем все конфигурации с истекшим сроком действия
	expiredConfigs, err := s.repo.GetExpiredConfigs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get expired configurations: %w", err)
	}

	// Обновляем статус на "expired" для каждой конфигурации
	for _, config := range expiredConfigs {
		err := s.UpdateStatus(ctx, config.ID, models.StatusExpired)
		if err != nil {
			return fmt.Errorf("failed to update status for config %d: %w", config.ID, err)
		}
	}

	return nil
}

func (s *ConfigurationService) GetExpiredConfigs(ctx context.Context) ([]models.Configuration, error) {
	return s.repo.GetExpiredConfigs(ctx)
}

func generateWireGuardKeyPair() (string, string, error) {
	// Генерация приватного ключа
	privateKeyCmd := exec.Command("wg", "genkey")
	privateKeyBytes, err := privateKeyCmd.Output()
	if err != nil {
		return "", "", err
	}
	privateKey := strings.TrimSpace(string(privateKeyBytes))

	// Генерация публичного ключа
	publicKeyCmd := exec.Command("wg", "pubkey")
	publicKeyCmd.Stdin = bytes.NewBufferString(privateKey)
	publicKeyBytes, err := publicKeyCmd.Output()
	if err != nil {
		return "", "", err
	}
	publicKey := strings.TrimSpace(string(publicKeyBytes))

	return privateKey, publicKey, nil
}
