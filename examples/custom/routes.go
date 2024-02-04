package main

import (
	"net/http"

	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/container"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/routing"
)

type Routes struct {
	routing.APIRoutes
}

func (r *Routes) Register() {
	r.AddRoute(http.MethodGet, "/message", container.Provide[IMessageController]().GetMessage)
}
