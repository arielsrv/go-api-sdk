package application

import (
	"context"

	"github.com/arielsrv/go-sdk-api/core/routing"
	"github.com/arielsrv/go-sdk-api/core/services"
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
