package container

import (
	"log"

	"go.uber.org/dig"
)

var Registry = dig.New()

func Provide[T any]() T {
	var instance T
	if err := Registry.Invoke(func(controller T) {
		instance = controller
	}); err != nil {
		log.Fatal(err)
	}

	return instance
}

func Inject(constructor interface{}, opts ...dig.ProvideOption) {
	if err := Registry.Provide(constructor, opts...); err != nil {
		log.Fatal(err)
	}
}
