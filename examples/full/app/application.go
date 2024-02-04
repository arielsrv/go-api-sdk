package app

import (
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/application"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/examples/full/app/workers"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/go-sdk-config/env"
	"log/slog"
	"os"
)

type Application struct {
	application.APIApplication
}

func (r *Application) Init() {
	// prometheus enabled (%err rate, apdex, etc.) see grafana (dev|uat|prod)
	r.UseMetrics()

	// optional, r.UseConfig() default  folder recommended src/config/**
	r.UseConfig(&env.Config{
		File: "example.yml", // recommended
		Path: "config",
		Logger: slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		})),
	})

	// optional, {protocol}://{host}:{port}/swagger/index.html
	r.UseSwagger()

	// optional
	r.UseLogger()

	// optional
	r.UseCors()

	// optional
	r.UseRequestID()

	// optional, recommended default folder src/app/views
	r.UseViews()

	// optionally, the app will be available after warmup
	r.RegisterWarmup(new(ApplicationWarmup))

	// optionally, the worker will be available after startup
	r.RegisterWorkers(new(workers.HelloWorldService))

	// recommended
	r.RegisterDependencyInjectionModule(new(ApplicationModule))

	// required
	r.RegisterRoutes(new(Routes))

	// required
	r.Build()
}
