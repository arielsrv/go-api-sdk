package application

import (
	"github.com/arielsrv/go-sdk-api/core/container"
	"github.com/arielsrv/go-sdk-api/core/routing"
	"github.com/arielsrv/go-sdk-api/core/services"
)

type IApplication interface {
	Init()
	RegisterServer(server Server)
}

type Application interface {
	RegisterWarmup(applicationWarmup Warmup)
	RegisterRoutes(router routing.Router)
	RegisterDependencyInjectionModule(applicationModule container.ApplicationModule)
	RegisterWorkers(workers ...services.IHostedService)
	Build()
}
