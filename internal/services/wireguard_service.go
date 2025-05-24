package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"wg_api/config"
	"wg_api/internal/models"
	"wg_api/pkg/shell"
)

type WireGuardService struct {
	executor               *shell.Executor
	configPath             string
	configService          *ConfigurationService
	wireguardContainerName string
}

func NewWireGuardService(
	executor *shell.Executor,
	cfg *config.Config,
	configService *ConfigurationService,
) *WireGuardService {
	return &WireGuardService{
		executor:               executor,
		configPath:             cfg.WireGuard.ConfigPath,
		configService:          configService,
		wireguardContainerName: cfg.WireGuard.ContainerName, // Имя контейнера из конфигурации
	}
}

func (s *WireGuardService) ApplyConfiguration(ctx context.Context, config *models.Configuration) error {
	// Создаем временный файл
	tempFile, err := os.CreateTemp("", "wg-config-*.conf")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	// Записываем конфиг
	if err := s.writeWireguardConfig(tempFile.Name(), config); err != nil {
		return err
	}

	// Копируем конфиг в общий том
	copyCmd := fmt.Sprintf("cp %s %s", tempFile.Name(), s.configPath)
	if _, err := s.executor.Execute("sh", "-c", copyCmd); err != nil {
		return err
	}

	// Перезагружаем WireGuard
	return s.reloadWireGuardInterface()
}

func (s *WireGuardService) writeWireguardConfig(tempPath string, config *models.Configuration) error {
	// Чтение текущего конфига
	currentConfig, err := os.ReadFile(s.configPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Добавляем нового пира
	newConfig := strings.Replace(string(currentConfig),
		"[Interface]",
		fmt.Sprintf("[Interface]\n\n%s", config.ToWireGuardConfig()),
		1)

	return os.WriteFile(tempPath, []byte(newConfig), 0644)
}

// Новый метод для перезагрузки интерфейса WireGuard в контейнере
func (s *WireGuardService) reloadWireGuardInterface() error {
	// Команда для остановки интерфейса в контейнере
	downCmd := fmt.Sprintf("docker exec %s wg-quick down wg0", s.wireguardContainerName)
	if _, err := s.executor.Execute("bash", "-c", downCmd); err != nil {
		return fmt.Errorf("failed to down wireguard interface in container: %w", err)
	}

	// Небольшая пауза для завершения процессов
	time.Sleep(1 * time.Second)

	// Команда для запуска интерфейса в контейнере
	upCmd := fmt.Sprintf("docker exec %s wg-quick up wg0", s.wireguardContainerName)
	if _, err := s.executor.Execute("bash", "-c", upCmd); err != nil {
		return fmt.Errorf("failed to up wireguard interface in container: %w", err)
	}

	return nil
}

func (s *WireGuardService) addOrUpdatePeer(peerConfig string, config *models.Configuration) error {
	// Проверяем, существует ли пир с таким публичным ключом
	existingPeer, err := s.getPeerByPublicKey(config.PublicKey)
	if err != nil {
		return err
	}

	if existingPeer != "" {
		// Удаляем существующий пир перед добавлением обновленного
		if err := s.removePeer(config.PublicKey, config); err != nil {
			return err
		}
	}

	// Добавляем новый пир в конфигурацию
	if err := s.addPeerToConfig(peerConfig, config); err != nil {
		return err
	}

	return nil
}

func (s *WireGuardService) getPeerByPublicKey(publicKey string) (string, error) {
	// Выполняем команду wg в контейнере
	cmd := fmt.Sprintf("docker exec %s wg show wg0 peers", s.wireguardContainerName)
	output, err := s.executor.Execute("bash", "-c", cmd)
	if err != nil {
		return "", fmt.Errorf("failed to get wireguard peers from container: %w", err)
	}

	peers := strings.Split(output, "\n")
	for _, peer := range peers {
		if strings.TrimSpace(peer) == publicKey {
			return peer, nil
		}
	}

	return "", nil
}

func (s *WireGuardService) removePeer(publicKey string, config *models.Configuration) error {
	// Чтение текущей конфигурации
	configData, err := s.readWireguardConfig()
	if err != nil {
		return err
	}

	// Поиск и удаление секции пира с указанным публичным ключом
	lines := strings.Split(configData, "\n")
	var newConfig []string
	skipLines := false

	for _, line := range lines {
		if strings.HasPrefix(line, "[Peer]") {
			skipLines = false // Сбрасываем флаг для новой секции пира
		}

		if strings.Contains(line, "PublicKey") && strings.Contains(line, publicKey) {
			skipLines = true // Начинаем пропускать строки для этого пира
			continue
		}

		if !skipLines {
			newConfig = append(newConfig, line)
		}
	}

	// Запись обновленной конфигурации
	return s.writeWireguardConfig(strings.Join(newConfig, "\n"), config)
}

func (s *WireGuardService) addPeerToConfig(peerConfig string, config *models.Configuration) error {
	configData, err := s.readWireguardConfig()
	if err != nil {
		return err
	}

	newConfig := configData + peerConfig
	return s.writeWireguardConfig(newConfig, config) // Добавляем второй аргумент
}

func (s *WireGuardService) readWireguardConfig() (string, error) {
	// Читаем конфигурацию из общего тома
	output, err := s.executor.Execute("cat", s.configPath)
	if err != nil {
		return "", fmt.Errorf("failed to read wireguard config: %w", err)
	}
	return output, nil
}

// func (s *WireGuardService) writeWireguardConfig(config string) error {
// 	// Записываем конфигурацию в общий том
// 	cmd := fmt.Sprintf("echo '%s' > %s", config, s.configPath)
// 	if _, err := s.executor.Execute("bash", "-c", cmd); err != nil {
// 		return fmt.Errorf("failed to write wireguard config: %w", err)
// 	}

// 	return nil
// }

func (s *WireGuardService) UpdatePeerHandshakes(ctx context.Context) error {
	// Получаем все активные конфигурации
	configs, err := s.configService.GetAllByStatus(ctx, models.StatusPaid)
	if err != nil {
		return fmt.Errorf("failed to get active configurations: %w", err)
	}

	// Получаем информацию о handshake из контейнера WireGuard
	cmd := fmt.Sprintf("docker exec %s wg show wg0 dump", s.wireguardContainerName)
	output, err := s.executor.Execute("bash", "-c", cmd)
	if err != nil {
		return fmt.Errorf("failed to get wireguard dump from container: %w", err)
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		publicKey := fields[1]
		lastHandshakeStr := fields[4]

		// Пропускаем, если нет handshake
		if lastHandshakeStr == "0" {
			continue
		}

		// Конвертируем Unix timestamp в time.Time
		lastHandshakeUnix, err := parseInt64(lastHandshakeStr)
		if err != nil {
			continue
		}

		lastHandshake := time.Unix(lastHandshakeUnix, 0)

		// Обновляем handshake для соответствующей конфигурации
		for _, config := range configs {
			if config.PublicKey == publicKey {
				err = s.configService.UpdateLatestHandshake(ctx, config.ID, lastHandshake)
				if err != nil {
					return fmt.Errorf("failed to update handshake for config %d: %w", config.ID, err)
				}
				break
			}
		}
	}

	return nil
}

func parseInt64(s string) (int64, error) {
	var i int64
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}
