package app

import (
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/examples/full/app/controllers"
	"net/http"

	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/routing"
)

// Routes is the routes for the application
type Routes struct {
	routing.APIRoutes
}

// Register and forward the request to the appropriate controller from DI Module
func (r *Routes) Register() {
	r.AddRoute(http.MethodGet, "/", routing.To[controllers.IHomeController]().Index)
	r.AddRoute(http.MethodGet, "/message", routing.To[controllers.IMessageController]().GetMessage)
}
