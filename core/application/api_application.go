package application

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/arielsrv/go-sdk-api/core/container"
	"github.com/arielsrv/go-sdk-api/core/routing"
	"github.com/arielsrv/go-sdk-api/core/services"
	log "gitlab.com/iskaypetcom/digital/sre/tools/dev/go-logger"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/go-sdk-config/env"
)

type APIApplication struct {
	server    Server
	appConfig AppConfig
	routes    []routing.Route
}

func (r *APIApplication) UseMetrics() {
	r.appConfig.Metrics = true
}

func (r *APIApplication) UseViews() {
	r.appConfig.Views = true
}

func (r *APIApplication) UseSwagger() {
	r.appConfig.Swagger = true
}

func (r *APIApplication) UseRequestID() {
	r.appConfig.RequestID = true
}

func (r *APIApplication) UseLogger() {
	r.appConfig.Logger = true
}

func (r *APIApplication) UseCors() {
	r.appConfig.Cors = true
}

func (r *APIApplication) UseRecovery() {
	r.appConfig.Recovery = true
}

func (r *APIApplication) UseStatic(prefix string, root string) {
	r.appConfig.Static = &routing.Static{
		Prefix: prefix,
		Root:   root,
	}
}

func (r *APIApplication) RegisterWarmup(applicationWarmup Warmup) {
	r.appConfig.ApplicationWarmup = applicationWarmup
}

func (r *APIApplication) RegisterWorkers(workers ...services.IHostedService) {
	r.appConfig.Workers = workers
}

func (r *APIApplication) RegisterRoutes(router routing.Router) {
	router.Register()
	r.routes = router.GetRoutes()
}

func (r *APIApplication) Build() {
	r.server.Configure(r.appConfig)
	r.server.RegisterRoutes(r.routes)
}

func (r *APIApplication) UseConfig(config ...*env.Config) {
	r.appConfig.Enabled = true
	if len(config) > 0 {
		env.SetConfigPath(config[0].Path)
		env.SetConfigFile(config[0].File)
		env.SetLogger(config[0].Logger)
	} else {
		env.SetConfigPath(filepath.Join("src", "resources", "config"))
		env.SetConfigFile("config.yml")
		env.SetLogger(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		})))
	}
	if err := env.Load(); err != nil {
		log.Error(err)
	}
}

func (r *APIApplication) RegisterServer(server Server) {
	r.server = server
}

func (r *APIApplication) RegisterDependencyInjectionModule(applicationModule container.ApplicationModule) {
	applicationModule.Configure()
}

type DefaultApplication struct {
	APIApplication
}

func (r *DefaultApplication) Init() {
	r.UseRecovery()
	r.UseLogger()
	r.UseMetrics()
	r.Build()
}
