package application

import (
	"github.com/arielsrv/go-sdk-api/core/routing"
	"github.com/arielsrv/go-sdk-api/core/services"
)

type AppConfig struct {
	Recovery          bool
	Swagger           bool
	RequestID         bool
	Logger            bool
	Cors              bool
	Metrics           bool
	Views             bool
	Enabled           bool
	ApplicationWarmup Warmup
	Workers           []services.IHostedService
	Static            *routing.Static
}
