package handlers

import (
	"context"
	"errors"
	"short_url/internal/handlers/middlewares"
	"short_url/internal/models"
	myLog "short_url/pkg/logger"

	"github.com/boombuler/barcode"
	"github.com/gin-gonic/gin"
)

// manageService Интерфейс к сервису, осуществляющему управление пользователя ссылками
type manageService interface {
	FindLink(ctx context.Context, link string) (models.LinkDataDTO, error)
	GetAllLinks(ctx context.Context, username string) ([]models.LinkDataDTO, error)
	DeleteLink(ctx context.Context, username, link string) error
	CreateLink(ctx context.Context, fullUrl, custom string, exp int, user models.JWTUserInfo) (models.LinkDataDTO, error)
	CreateQR(ctx context.Context, url, link string) (barcode.Barcode, error)
}

// ManageHandlerConfig Конфигурация для ManageHandler
type ManageHandlerConfig struct {
	Router        *gin.Engine
	ManageService manageService
	Middlware		*middlewares.Middlewares
	Logger			*myLog.Log
}

// ManageHandler Для логирования и регистрации хендлеров
type ManageHandler struct {
	manageService manageService
	middleware		*middlewares.Middlewares
	logger			*myLog.Log
}

// GetUserInfo Возвращает информацию о пользователе и его подписке из контекста
func (h *ManageHandler) GetUserInfo(ctx *gin.Context) models.JWTUserInfo {

	// Получаем данные из контекста
	info, ok := ctx.Get(middlewares.UserInfo)
	if !ok {
		h.logger.Error("failed to get user info from ctx")

		return models.JWTUserInfo{
			Err:	errors.New("failed to get user info"),
		}
	}

	// Преобразуем их в структуру для ответа
	user, ok := info.(models.JWTUserInfo)
	if !ok {
		h.logger.Error("failed to get JWT user info")

		return models.JWTUserInfo{
			Err:	errors.New("failed to JWT user info"),
		}
	}

	return user
}

// RegisterManageHandler Фабрика для ManageHandler
func RegisterManageHandler(c *ManageHandlerConfig) {
	manageHandler := ManageHandler{
		manageService: c.ManageService,
		middleware:		c.Middlware,
		logger:			c.Logger,
	}

	g := c.Router.Group("v1")
	g.POST("/newlink",c.Middlware.AuthUser, manageHandler.CreateLink)
	g.DELETE("/links/:link",c.Middlware.AuthUser, manageHandler.DeleteLink)
	g.GET("/links",c.Middlware.AuthUser, manageHandler.GetAllLinks)
	g.GET("/:link", manageHandler.LinkRedirect)
	g.GET("/links/qr/:link",c.Middlware.AuthUser, manageHandler.CreateCode)
	g.GET("/links/:link",c.Middlware.AuthUser, manageHandler.GetLink)
}