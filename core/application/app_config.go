package application

import (
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/routing"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/services"
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
