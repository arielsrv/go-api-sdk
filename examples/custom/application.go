package main

import (
	"github.com/arielsrv/go-sdk-api/core/application"
)

type Application struct {
	application.APIApplication
}

func (r *Application) Init() {
	r.UseMetrics()
	r.UseSwagger()
	r.RegisterDependencyInjectionModule(new(ApplicationModule))
	r.RegisterRoutes(new(Routes))
	r.Build()
}
