package application

import (
	"context"

	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/routing"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/services"
)

type Server interface {
	Start()
	Shutdown() error
	Join() error
	RegisterRoutes(routes []routing.Route)
	Configure(config AppConfig)
	On(application IApplication)
}

type SelfHosting interface {
	AddHostedService(hostedService services.IHostedService)
	GetHostedServices() []services.IHostedService
}

type Warmupper interface {
	Loaded(ctx context.Context) bool
	SetIsReady(value bool)
	HasWarmup() bool
}
