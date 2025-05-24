package scheduler

import (
	"context"
	"log"
	"time"
	"wg_api/internal/models"
	"wg_api/internal/services"
)

type Scheduler struct {
	configService    *services.ConfigurationService
	wireguardService *services.WireGuardService
	interval         time.Duration
	stop             chan struct{}
}

func NewScheduler(
	configService *services.ConfigurationService,
	wireguardService *services.WireGuardService,
	interval time.Duration,
) *Scheduler {
	return &Scheduler{
		configService:    configService,
		wireguardService: wireguardService,
		interval:         interval,
		stop:             make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.runTasks()
			case <-s.stop:
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stop)
}

func (s *Scheduler) runTasks() {
	ctx := context.Background()

	// Задача 1: Обновление статусов конфигураций
	if err := s.configService.UpdateExpiredConfigurations(ctx); err != nil {
		log.Printf("Error updating expired configurations: %v", err)
	}

	// Задача 2: Обновление информации о handshake
	if err := s.wireguardService.UpdatePeerHandshakes(ctx); err != nil {
		log.Printf("Error updating peer handshakes: %v", err)
	}

	log.Println("Scheduled tasks completed successfully")
}

func (s *Scheduler) updateExpiredConfigurations(ctx context.Context) error {
	now := time.Now()

	// Конфигурации с истекшим сроком
	expiredConfigs, err := s.configService.GetExpiredConfigs(ctx)
	if err != nil {
		return err
	}

	for i := range expiredConfigs {
		config := &expiredConfigs[i] // Работаем с указателем

		if config.Status == models.StatusPaid {
			config.Status = models.StatusExpired
			config.ExpirationTime = now.AddDate(0, 0, 7)
			if err := s.configService.Update(ctx, config); err != nil {
				return err
			}
		}

		if config.Status == models.StatusExpired && now.After(config.ExpirationTime) {
			if err := s.configService.Delete(ctx, config.ID); err != nil {
				return err
			}
			if err := s.wireguardService.ApplyConfiguration(ctx, config); err != nil {
				return err
			}
		}
	}

	return nil
}
