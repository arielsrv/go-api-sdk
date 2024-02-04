package container_test

import (
	"testing"

	"github.com/arielsrv/go-sdk-api/core/container"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"
)

func TestDependencyInjectionModule_Bind(t *testing.T) {
	module := &container.DependencyInjectionModule{}
	module.Bind(NewTestObject, dig.As(new(ITestObject)))

	actual := container.Provide[ITestObject]()
	require.NotNil(t, actual)
	require.Equal(t, "test", actual.Test())
}

type ITestObject interface {
	Test() string
}

type TestObject struct {
}

func NewTestObject() *TestObject {
	return &TestObject{}
}

func (t *TestObject) Test() string {
	return "test"
}
