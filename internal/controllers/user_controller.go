package controllers

import (
	"net/http"
	"strconv"
	"wg_api/internal/models"
	"wg_api/internal/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *services.UserService
}

// NewUserController создает новый экземпляр контроллера пользователей
func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}

// Create godoc
// @Summary Создание пользователя
// @Description Создает нового пользователя в системе
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "Данные пользователя"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string "Ошибка валидации"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /users [post]
func (c *UserController) Create(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Create(ctx, &user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// GetByID godoc
// @Summary Получение пользователя по ID
// @Description Возвращает информацию о пользователе по его ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "Некорректный ID пользователя"
// @Failure 404 {object} map[string]string "Пользователь не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /users/{id} [get]
func (c *UserController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := c.service.GetByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// GetAll godoc
// @Summary Получение списка всех пользователей
// @Description Возвращает список всех пользователей в системе
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /users [get]
func (c *UserController) GetAll(ctx *gin.Context) {
	users, err := c.service.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// Update godoc
// @Summary Обновление пользователя
// @Description Обновляет данные существующего пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param user body models.User true "Обновленные данные пользователя"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "Некорректный ID пользователя или данные"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /users/{id} [put]
func (c *UserController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = uint(id)
	if err := c.service.Update(ctx, &user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// Delete godoc
// @Summary Удаление пользователя
// @Description Удаляет пользователя из системы
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} map[string]string "Сообщение об успешном удалении"
// @Failure 400 {object} map[string]string "Некорректный ID пользователя"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /users/{id} [delete]
func (c *UserController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := c.service.Delete(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// RegisterRoutes регистрирует маршруты для API пользователей
func (c *UserController) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.POST("", c.Create)
		users.GET("", c.GetAll)
		users.GET("/:id", c.GetByID)
		users.PUT("/:id", c.Update)
		users.DELETE("/:id", c.Delete)
	}
}
