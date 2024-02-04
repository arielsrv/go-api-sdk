package app

import (
	"github.com/arielsrv/go-sdk-api/core/container"
	"github.com/arielsrv/go-sdk-api/examples/full/app/controllers"
	"github.com/arielsrv/go-sdk-api/examples/full/app/services"
	"go.uber.org/dig"
)

type ApplicationModule struct {
	container.DependencyInjectionModule
}

// Configure the application module
func (r *ApplicationModule) Configure() {
	r.Bind(services.NewMessageService, dig.As(new(services.IMessageService)))
	r.Bind(controllers.NewMessageController, dig.As(new(controllers.IMessageController)))
	r.Bind(controllers.NewHomeController, dig.As(new(controllers.IHomeController)))
}
