package container_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/container"
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
