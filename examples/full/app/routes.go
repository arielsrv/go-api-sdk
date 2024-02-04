package app

import (
	"github.com/arielsrv/go-sdk-api/examples/full/app/controllers"
	"net/http"

	"github.com/arielsrv/go-sdk-api/core/routing"
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
