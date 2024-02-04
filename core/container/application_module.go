package container

import (
	"go.uber.org/dig"
)

type ApplicationModule interface {
	Configure()
}

type InjectionModule interface {
	Bind(constructor interface{}, opts ...dig.ProvideOption) DependencyInjectionModule
}

type DependencyInjectionModule struct {
}

func (r *DependencyInjectionModule) Bind(constructor interface{}, opts ...dig.ProvideOption) DependencyInjectionModule {
	Inject(constructor, opts...)
	return *r
}
