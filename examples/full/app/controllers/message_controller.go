package controllers

import (
	"github.com/arielsrv/go-sdk-api/core/routing"
	"github.com/arielsrv/go-sdk-api/examples/full/app/services"
)

type IMessageController interface {
	GetMessage(ctx *routing.HTTPContext) error
}

type MessageController struct {
	messageService services.IMessageService
}

func NewMessageController(messageService services.IMessageService) *MessageController {
	return &MessageController{
		messageService: messageService,
	}
}

// GetMessage godoc
// @Summary	Get message
// @Description	Get message
// @Tags	Messages
// @Success	200
// @Produce	json
// @Produce	application/problem+json
// @Success	200	{object} string
// @Failure 404 {object} errorx.Error
// @Failure 500 {object} errorx.Error
// @Router	/message [get].
func (r *MessageController) GetMessage(ctx *routing.HTTPContext) error {
	message := r.messageService.GetMessage()
	return ctx.SendString(message)
}
