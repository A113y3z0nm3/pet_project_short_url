package handlers

import (
	"net/http"
	"short_url/internal/models"

	"github.com/gin-gonic/gin"
)

// QiwiSub Создает счет в QIWI P2P для покупки подписки
func (h *PayHandler) QiwiSub(ctx *gin.Context) {
	// Получаем информацию о пользователе
	user, err := h.middleware.GetUserInfo(ctx)
	if err != nil {
		//
	}

	// Если пользователь подписчик, закрываем ручку и отдаем ответ
	if user.Subscribe == models.Sub {
		ctx.JSON(http.StatusOK, gin.H{
			"error": "the user already has a subscribe",
		})
	}

	// Получаем желаемый срок подписки из path
	sub, err := getSubFromParam(ctx, h.prices)
	if err != nil {
		//
	}

	// Создаем счет на оплату и получаем ссылку на него
	result, err := h.qiwiService.BillRequest(ctx, sub, user.Username)
	if err != nil {
		//
	}

	// Отдаем ссылку на счет клиенту
	ctx.JSON(http.StatusOK, gin.H{
		"payUrl": result,
	})

	return
}