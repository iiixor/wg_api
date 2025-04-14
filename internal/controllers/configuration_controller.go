package controllers

import (
	"net/http"
	"strconv"
	"wg_api/internal/models"
	"wg_api/internal/services"

	"github.com/gin-gonic/gin"
)

type ConfigurationController struct {
	service          *services.ConfigurationService
	wireguardService *services.WireGuardService
}

// NewConfigurationController создает новый экземпляр контроллера конфигураций
func NewConfigurationController(
	service *services.ConfigurationService,
	wireguardService *services.WireGuardService,
) *ConfigurationController {
	return &ConfigurationController{
		service:          service,
		wireguardService: wireguardService,
	}
}

// Create godoc
// @Summary Создание конфигурации
// @Description Создает новую конфигурацию WireGuard
// @Tags configurations
// @Accept json
// @Produce json
// @Param config body models.Configuration true "Данные конфигурации"
// @Success 201 {object} models.Configuration
// @Failure 400 {object} map[string]string "Ошибка валидации"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /configurations [post]
func (c *ConfigurationController) Create(ctx *gin.Context) {
	var config models.Configuration
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Create(ctx, &config); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Если конфигурация оплачена, применяем изменения в WireGuard
	if config.Status == models.StatusPaid {
		if err := c.wireguardService.ApplyConfiguration(ctx, &config); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "config created but failed to apply wireguard changes: " + err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusCreated, config)
}

// GetByID godoc
// @Summary Получение конфигурации по ID
// @Description Возвращает информацию о конфигурации WireGuard по её ID
// @Tags configurations
// @Accept json
// @Produce json
// @Param id path int true "ID конфигурации"
// @Success 200 {object} models.Configuration
// @Failure 400 {object} map[string]string "Некорректный ID конфигурации"
// @Failure 404 {object} map[string]string "Конфигурация не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /configurations/{id} [get]
func (c *ConfigurationController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid configuration ID"})
		return
	}

	config, err := c.service.GetByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, config)
}

// GetAll godoc
// @Summary Получение всех конфигураций
// @Description Возвращает список всех конфигураций WireGuard
// @Tags configurations
// @Accept json
// @Produce json
// @Success 200 {array} models.Configuration
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /configurations [get]
func (c *ConfigurationController) GetAll(ctx *gin.Context) {
	configs, err := c.service.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, configs)
}

// GetByUserID godoc
// @Summary Получение конфигураций пользователя
// @Description Возвращает список конфигураций WireGuard для конкретного пользователя
// @Tags configurations
// @Accept json
// @Produce json
// @Param userId path int true "ID пользователя"
// @Success 200 {array} models.Configuration
// @Failure 400 {object} map[string]string "Некорректный ID пользователя"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /configurations/user/{userId} [get]
func (c *ConfigurationController) GetByUserID(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	configs, err := c.service.GetByUserID(ctx, uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, configs)
}

// Update godoc
// @Summary Обновление конфигурации
// @Description Обновляет существующую конфигурацию WireGuard
// @Tags configurations
// @Accept json
// @Produce json
// @Param id path int true "ID конфигурации"
// @Param config body models.Configuration true "Обновленные данные конфигурации"
// @Success 200 {object} models.Configuration
// @Failure 400 {object} map[string]string "Некорректный ID конфигурации или данные"
// @Failure 404 {object} map[string]string "Конфигурация не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /configurations/{id} [put]
func (c *ConfigurationController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid configuration ID"})
		return
	}

	// Получаем текущую конфигурацию для проверки изменения статуса
	oldConfig, err := c.service.GetByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var config models.Configuration
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.ID = uint(id)
	if err := c.service.Update(ctx, &config); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Если статус изменился, применяем изменения в WireGuard
	if oldConfig.Status != config.Status {
		if err := c.wireguardService.ApplyConfiguration(ctx, &config); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "config updated but failed to apply wireguard changes: " + err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, config)
}

// UpdateStatus godoc
// @Summary Обновление статуса конфигурации
// @Description Обновляет статус конфигурации WireGuard
// @Tags configurations
// @Accept json
// @Produce json
// @Param id path int true "ID конфигурации"
// @Param status body object true "Новый статус" example={"status":"paid"}
// @Success 200 {object} map[string]string "Сообщение об успешном обновлении"
// @Failure 400 {object} map[string]string "Некорректный ID конфигурации или статус"
// @Failure 404 {object} map[string]string "Конфигурация не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /configurations/{id}/status [patch]
func (c *ConfigurationController) UpdateStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid configuration ID"})
		return
	}

	type StatusUpdate struct {
		Status models.ConfigStatus `json:"status" binding:"required"`
	}

	var update StatusUpdate
	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UpdateStatus(ctx, uint(id), update.Status); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Получаем обновленную конфигурацию и применяем изменения в WireGuard
	config, err := c.service.GetByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := c.wireguardService.ApplyConfiguration(ctx, config); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "status updated but failed to apply wireguard changes: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "status updated successfully"})
}

// Delete godoc
// @Summary Удаление конфигурации
// @Description Удаляет конфигурацию WireGuard
// @Tags configurations
// @Accept json
// @Produce json
// @Param id path int true "ID конфигурации"
// @Success 200 {object} map[string]string "Сообщение об успешном удалении"
// @Failure 400 {object} map[string]string "Некорректный ID конфигурации"
// @Failure 404 {object} map[string]string "Конфигурация не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /configurations/{id} [delete]
func (c *ConfigurationController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid configuration ID"})
		return
	}

	// Получаем конфигурацию перед удалением
	config, err := c.service.GetByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем статус на "deletion" и применяем изменения в WireGuard
	config.Status = models.StatusDeletion
	if err := c.wireguardService.ApplyConfiguration(ctx, config); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to remove from wireguard before deletion: " + err.Error(),
		})
		return
	}

	if err := c.service.Delete(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "configuration deleted successfully"})
}

// RegisterRoutes регистрирует маршруты для API конфигураций
func (c *ConfigurationController) RegisterRoutes(router *gin.RouterGroup) {
	configs := router.Group("/configurations")
	{
		configs.POST("", c.Create)
		configs.GET("", c.GetAll)
		configs.GET("/:id", c.GetByID)
		configs.GET("/user/:userId", c.GetByUserID)
		configs.PUT("/:id", c.Update)
		configs.PATCH("/:id/status", c.UpdateStatus)
		configs.DELETE("/:id", c.Delete)
	}
}
