package application

import (
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/container"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/routing"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/services"
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
