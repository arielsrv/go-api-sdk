package main

import (
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/container"
	"go.uber.org/dig"
)

type ApplicationModule struct {
	container.DependencyInjectionModule
}

func (r *ApplicationModule) Configure() {
	r.Bind(NewMessageController, dig.As(new(IMessageController)))
}
