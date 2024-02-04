package routing_test

import (
	"net/http"
	"testing"

	"github.com/arielsrv/go-sdk-api/core/container"
	"github.com/arielsrv/go-sdk-api/core/routing"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

func TestAPIRoutes_GetRoutes(t *testing.T) {
	apiRoutes := new(routing.APIRoutes)
	apiRoutes.AddRoute(http.MethodGet, "/get", func(ctx *routing.HTTPContext) error {
		return ctx.SendString("Hello, World!")
	})

	apiRoutes.AddRoute(http.MethodPatch, "/patch", func(ctx *routing.HTTPContext) error {
		return ctx.SendString("Hello, World!")
	})

	actual := apiRoutes.GetRoutes()

	assert.Len(t, actual, 2)

	assert.Equal(t, http.MethodGet, actual[0].Method)
	assert.Equal(t, "/get", actual[0].Path)

	assert.Equal(t, http.MethodPatch, actual[1].Method)
	assert.Equal(t, "/patch", actual[1].Path)
}

func TestAPIRoutes_GetStatics(t *testing.T) {
	apiRoutes := new(routing.APIRoutes)
	apiRoutes.AddStatic("/static", "/static")

	actual := apiRoutes.GetStatics()

	assert.Len(t, actual, 1)

	assert.Equal(t, "/static", actual[0].Prefix)
	assert.Equal(t, "/static", actual[0].Root)
}

type Command interface {
	Execute() string
}

type MessageCommand struct {
}

func (c *MessageCommand) Execute() string {
	return "Hello World!"
}

func TestTo(t *testing.T) {
	container.Inject(func() Command {
		return &MessageCommand{}
	}, dig.As(new(Command)))

	actual := routing.To[Command]()

	assert.NotNil(t, actual)
	assert.Equal(t, "Hello World!", actual.Execute())
}
