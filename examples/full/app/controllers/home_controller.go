package controllers

import (
	"github.com/arielsrv/go-sdk-api/core/routing"
)

type IHomeController interface {
	Index(ctx *routing.HTTPContext) error
}

type HomeController struct {
}

func NewHomeController() *HomeController {
	return &HomeController{}
}

func (h HomeController) Index(ctx *routing.HTTPContext) error {
	return ctx.Render("home/index", map[string]string{
		"title": "Home",
		"body":  "This is the home page",
	})
}
