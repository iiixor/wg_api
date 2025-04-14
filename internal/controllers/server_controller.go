package controllers

import (
	"net/http"
	"strconv"
	"wg_api/internal/models"
	"wg_api/internal/services"

	"github.com/gin-gonic/gin"
)

type ServerController struct {
	service *services.ServerService
}

// NewServerController создает новый экземпляр контроллера серверов
func NewServerController(service *services.ServerService) *ServerController {
	return &ServerController{service: service}
}

// Create godoc
// @Summary Создание сервера
// @Description Создает новый WireGuard сервер в системе
// @Tags servers
// @Accept json
// @Produce json
// @Param server body models.Server true "Данные сервера"
// @Success 201 {object} models.Server
// @Failure 400 {object} map[string]string "Ошибка валидации"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /servers [post]
func (c *ServerController) Create(ctx *gin.Context) {
	var server models.Server
	if err := ctx.ShouldBindJSON(&server); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Create(ctx, &server); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, server)
}

// GetByID godoc
// @Summary Получение сервера по ID
// @Description Возвращает информацию о WireGuard сервере по его ID
// @Tags servers
// @Accept json
// @Produce json
// @Param id path int true "ID сервера"
// @Success 200 {object} models.Server
// @Failure 400 {object} map[string]string "Некорректный ID сервера"
// @Failure 404 {object} map[string]string "Сервер не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /servers/{id} [get]
func (c *ServerController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid server ID"})
		return
	}

	server, err := c.service.GetByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, server)
}

// GetAll godoc
// @Summary Получение списка всех серверов
// @Description Возвращает список всех WireGuard серверов в системе
// @Tags servers
// @Accept json
// @Produce json
// @Success 200 {array} models.Server
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /servers [get]
func (c *ServerController) GetAll(ctx *gin.Context) {
	servers, err := c.service.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, servers)
}

// Update godoc
// @Summary Обновление сервера
// @Description Обновляет данные существующего WireGuard сервера
// @Tags servers
// @Accept json
// @Produce json
// @Param id path int true "ID сервера"
// @Param server body models.Server true "Обновленные данные сервера"
// @Success 200 {object} models.Server
// @Failure 400 {object} map[string]string "Некорректный ID сервера или данные"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /servers/{id} [put]
func (c *ServerController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid server ID"})
		return
	}

	var server models.Server
	if err := ctx.ShouldBindJSON(&server); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	server.ID = uint(id)
	if err := c.service.Update(ctx, &server); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, server)
}

// Delete godoc
// @Summary Удаление сервера
// @Description Удаляет WireGuard сервер из системы
// @Tags servers
// @Accept json
// @Produce json
// @Param id path int true "ID сервера"
// @Success 200 {object} map[string]string "Сообщение об успешном удалении"
// @Failure 400 {object} map[string]string "Некорректный ID сервера"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /servers/{id} [delete]
func (c *ServerController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid server ID"})
		return
	}

	if err := c.service.Delete(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "server deleted successfully"})
}

// RegisterRoutes регистрирует маршруты для API серверов
func (c *ServerController) RegisterRoutes(router *gin.RouterGroup) {
	servers := router.Group("/servers")
	{
		servers.POST("", c.Create)
		servers.GET("", c.GetAll)
		servers.GET("/:id", c.GetByID)
		servers.PUT("/:id", c.Update)
		servers.DELETE("/:id", c.Delete)
	}
}
