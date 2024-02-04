package core_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/application"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/container"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/mocks"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/routing"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/go-restclient/rest"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/go-sdk-config/env"
	"go.uber.org/dig"
)

type ServerSuite struct {
	suite.Suite
	port   int
	server *core.FastHTTPServer
}

func (r *ServerSuite) SetupTest() {
	listener, err := net.Listen("tcp", ":0")
	r.Require().NoError(err)

	addr, ok := listener.Addr().(*net.TCPAddr)
	r.True(ok)

	port := addr.Port
	err = listener.Close()
	r.Require().NoError(err)

	r.port = port
	r.server = core.NewServer(port)
}

func (r *ServerSuite) TearDownTest() {
	r.Require().NoError(r.server.Join())
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}

func (r *ServerSuite) TestServerStart() {
	r.server.On(new(Application))
	r.server.Start()

	go func() {
		time.Sleep(200 * time.Millisecond)

		rb := rest.RequestBuilder{
			BaseURL: fmt.Sprintf("http://0.0.0.0:%d", r.port),
		}

		response := rb.Get("/ping")
		r.Require().NoError(response.Err)
		r.Equal(http.StatusServiceUnavailable, response.StatusCode)
		r.Equal(http.StatusText(http.StatusServiceUnavailable), response.String())

		time.Sleep(200 * time.Millisecond)

		response = rb.Get("/ping")
		r.Require().NoError(response.Err)
		r.Equal(http.StatusOK, response.StatusCode)
		r.Equal(http.StatusText(http.StatusOK), response.String())

		response = rb.Get("/users")
		r.Require().NoError(response.Err)
		r.Equal(http.StatusOK, response.StatusCode)

		response = rb.Get("/orders")
		r.Require().NoError(response.Err)
		r.Equal(http.StatusOK, response.StatusCode)

		response = rb.Get("/metrics")
		r.Require().NoError(response.Err)
		r.Equal(http.StatusOK, response.StatusCode)

		response = rb.Get("/swagger/index.html")
		r.Require().NoError(response.Err)
		r.Equal(http.StatusOK, response.StatusCode)

		response = rb.Get("/swagger.yaml")
		r.Require().NoError(response.Err)
		r.Equal(http.StatusOK, response.StatusCode)

		r.Require().NoError(r.server.Shutdown())
	}()
}

func (r *ServerSuite) TestServerStart_NonWarmup() {
	r.server.On(new(NonWarmupApplication))
	r.server.AddHostedService(new(mocks.DummySQSConsumerService))
	r.server.Start()

	go func() {
		time.Sleep(200 * time.Millisecond)

		rb := rest.RequestBuilder{
			BaseURL: fmt.Sprintf("http://0.0.0.0:%d", r.port),
		}

		response := rb.Get("/ping")
		r.Require().NoError(response.Err)
		r.Equal(http.StatusOK, response.StatusCode)
		r.Equal(http.StatusText(http.StatusOK), response.String())

		response = rb.Get("/users")
		r.Require().NoError(response.Err)
		r.Equal(http.StatusOK, response.StatusCode)

		response = rb.Get("/orders")
		r.Require().NoError(response.Err)
		r.Equal(http.StatusOK, response.StatusCode)

		r.Require().NoError(r.server.Shutdown())
	}()
}

func (r *ServerSuite) TestServerStart_Default() {
	r.server.AddHostedService(new(mocks.DummySQSConsumerService))
	r.server.Start()

	go func() {
		time.Sleep(200 * time.Millisecond)

		rb := rest.RequestBuilder{
			BaseURL: fmt.Sprintf("http://0.0.0.0:%d", r.port),
		}

		response := rb.Get("/ping")
		r.Require().NoError(response.Err)
		r.Equal(http.StatusOK, response.StatusCode)
		r.Equal(http.StatusText(http.StatusOK), response.String())

		r.Require().NoError(r.server.Shutdown())
	}()
}

func (r *ServerSuite) TestServerStart_Join() {
	r.server.AddHostedService(new(mocks.DummySQSConsumerService))
	r.server.Start()

	go func() {
		time.Sleep(2000 * time.Millisecond)

		rb := rest.RequestBuilder{
			BaseURL: fmt.Sprintf("http://0.0.0.0:%d", r.port),
		}

		response := rb.Get("/ping")
		r.Require().NoError(response.Err)
		r.Equal(http.StatusOK, response.StatusCode)
		r.Equal(http.StatusText(http.StatusOK), response.String())

		r.Require().NoError(r.server.Shutdown())
	}()

	r.Require().NoError(r.server.Join())
}

type Application struct {
	application.APIApplication
}

func (r *Application) Init() {
	r.UseSwagger()
	r.UseRecovery()
	r.UseRequestID()
	r.UseLogger()
	r.UseCors()
	r.UseMetrics()
	r.UseConfig(&env.Config{
		File: "example.yml",
		Path: "config",
	})
	r.UseViews()
	r.UseStatic("/", "./../docs")
	r.RegisterWorkers(new(mocks.DummySQSConsumerService))
	r.RegisterDependencyInjectionModule(new(ApplicationModule))
	r.RegisterWarmup(new(ApplicationWarmup))
	r.RegisterRoutes(new(Routes))
	r.Build()
}

type NonWarmupApplication struct {
	application.APIApplication
}

func (r *NonWarmupApplication) Init() {
	r.UseRecovery()
	r.UseConfig()
	r.RegisterRoutes(new(Routes))
	r.Build()
}

type ApplicationModule struct {
	container.DependencyInjectionModule
}

func (r *ApplicationModule) Configure() {
	r.Bind(NewTestController, dig.As(new(ITestController)))
}

type Routes struct {
	routing.APIRoutes
}

type ApplicationWarmup struct {
}

func (r *ApplicationWarmup) Loaded(_ context.Context) bool {
	time.Sleep(time.Duration(400) * time.Millisecond)
	return true
}

type ITestController interface {
	GetUsers(ctx *routing.HTTPContext) error
	GetOrders(ctx *routing.HTTPContext) error
}

type TestController struct {
}

func (r TestController) GetUsers(ctx *routing.HTTPContext) error {
	return ctx.SendString("users")
}

func (r TestController) GetOrders(ctx *routing.HTTPContext) error {
	return ctx.SendString("orders")
}

func NewTestController() *TestController {
	return &TestController{}
}

func (r *Routes) Register() {
	r.AddRoute(http.MethodGet, "/users", container.Provide[ITestController]().GetUsers)
	r.AddRoute(http.MethodGet, "/orders", container.Provide[ITestController]().GetOrders)
}

type Route struct {
	Path string `validate:"required"`
}

func TestGuard(t *testing.T) {
	route := new(Route)
	err := core.Guard.Var(route.Path, "required")
	require.Error(t, err)

	route.Path = "/users"
	err = core.Guard.Var(route.Path, "required")
	require.NoError(t, err)
}
