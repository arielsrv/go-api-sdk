package app

import (
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/container"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/examples/full/app/controllers"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/examples/full/app/services"
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
