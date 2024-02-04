package controllers

import (
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/routing"
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
