package main

import (
	"github.com/arielsrv/go-sdk-api/core/container"
	"go.uber.org/dig"
)

type ApplicationModule struct {
	container.DependencyInjectionModule
}

func (r *ApplicationModule) Configure() {
	r.Bind(NewMessageController, dig.As(new(IMessageController)))
}
