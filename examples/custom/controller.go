package main

import "github.com/arielsrv/go-sdk-api/core/routing"

type IMessageController interface {
	GetMessage(ctx *routing.HTTPContext) error
}

type MessageController struct {
}

func NewMessageController() *MessageController {
	return &MessageController{}
}

func (r *MessageController) GetMessage(ctx *routing.HTTPContext) error {
	return ctx.SendString("Hello World")
}
