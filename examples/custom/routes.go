package main

import (
	"net/http"

	"github.com/arielsrv/go-sdk-api/core/container"
	"github.com/arielsrv/go-sdk-api/core/routing"
)

type Routes struct {
	routing.APIRoutes
}

func (r *Routes) Register() {
	r.AddRoute(http.MethodGet, "/message", container.Provide[IMessageController]().GetMessage)
}
