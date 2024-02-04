package container_test

import (
	"testing"

	"github.com/arielsrv/go-sdk-api/core/container"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

type Command interface {
	Execute() string
}

type MessageCommand struct {
}

func (c *MessageCommand) Execute() string {
	return "Hello World!"
}

func TestInject(t *testing.T) {
	container.Inject(func() Command {
		return &MessageCommand{}
	}, dig.As(new(Command)))

	actual := container.Provide[Command]()

	assert.NotNil(t, actual)
	assert.Equal(t, "Hello World!", actual.Execute())
}
