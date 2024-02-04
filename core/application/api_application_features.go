package application

import "gitlab.com/iskaypetcom/digital/sre/tools/dev/go-sdk-config/env"

type APIApplicationFeatures interface {
	UseMetrics()
	UseViews()
	UseSwagger()
	UseRequestID()
	UseLogger()
	UseCors()
	UseRecovery()
	UseConfig(config ...*env.Config)
	UseStatic(prefix string, root string)
}
