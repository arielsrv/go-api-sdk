package core

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/arielsrv/go-sdk-api/core/application"
	"github.com/arielsrv/go-sdk-api/core/collector"
	"github.com/arielsrv/go-sdk-api/core/errorx"
	"github.com/arielsrv/go-sdk-api/core/routing"
	"github.com/arielsrv/go-sdk-api/core/services"
	"github.com/arielsrv/go-sdk-api/core/subscriptions"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/html/v2"
	log "gitlab.com/iskaypetcom/digital/sre/tools/dev/go-logger"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/go-sdk-config/env"
	_ "go.uber.org/automaxprocs"
)

var _ = ""
var Guard = validator.New(validator.WithRequiredStructEnabled())

type FastHTTPServer struct {
	*fiber.App

	application    application.IApplication
	appConfig      application.AppConfig
	hostedServices []services.IHostedService
	routes         []routing.Route
	host           string
	port           int
	appName        string

	errChan                chan error
	onListen               *subscriptions.Notifier
	onListenSubscribers    []subscriptions.Listener
	applicationInitialized bool

	wg  sync.WaitGroup
	rw  sync.RWMutex
	mtx sync.Mutex
}

func (r *FastHTTPServer) Loaded(ctx context.Context) bool {
	return r.appConfig.ApplicationWarmup.Loaded(ctx)
}

func (r *FastHTTPServer) SetIsReady(value bool) {
	r.setIsReady(value)
}

func (r *FastHTTPServer) HasWarmup() bool {
	return r.appConfig.ApplicationWarmup != nil
}

func NewServer(port ...int) *FastHTTPServer {
	server := &FastHTTPServer{
		wg:                     sync.WaitGroup{},
		mtx:                    sync.Mutex{},
		rw:                     sync.RWMutex{},
		errChan:                make(chan error, 1),
		hostedServices:         make([]services.IHostedService, 0),
		routes:                 make([]routing.Route, 0),
		onListen:               subscriptions.NewNotifier(),
		onListenSubscribers:    make([]subscriptions.Listener, 0),
		host:                   "0.0.0.0",
		port:                   8081,
		applicationInitialized: false,
		appName:                "backend-sdk-api",
	}

	if len(port) > 0 {
		server.port = port[0]
	}

	return server
}

func (r *FastHTTPServer) Start() {
	ctx := context.Background()

	if r.application == nil {
		log.Warn("Application is not initialized, using default ...")
		r.On(new(application.DefaultApplication))
	}

	r.application.Init()

	r.wg.Add(1)
	go func(r *FastHTTPServer) {
		defer r.wg.Done()
		r.registerSubscribers(ctx)
		if r.appConfig.Enabled {
			if serverHost := env.Get("server.host"); !env.IsEmptyString(serverHost) {
				r.host = serverHost
			}
			r.port = env.GetInt("server.port", r.port)
		}
		addr := net.JoinHostPort(r.host, strconv.Itoa(r.port))
		r.App.Hooks().OnListen(func(addr fiber.ListenData) error {
			log.Println("server listening on", net.JoinHostPort(addr.Host, addr.Port))
			return r.onListen.Send(true)
		})
		r.errChan <- r.App.Listen(addr)
	}(r)

	go func() {
		r.wg.Wait()
		close(r.errChan)
		r.onListen.Close()
	}()
}

func (r *FastHTTPServer) Configure(appConfig application.AppConfig) {
	r.appName = env.GetString("APP_NAME", r.appName)

	config := fiber.Config{
		AppName:               r.appName,
		ServerHeader:          "x-backend-sdk-api",
		DisableStartupMessage: true,
		EnablePrintRoutes:     false,
		ErrorHandler:          errorx.OnErrorHandler,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	}

	if appConfig.Views {
		viewsFolder := env.GetString("views.folder", "src/app/views")
		config.Views = html.New(viewsFolder, ".html")
	}

	r.App = fiber.New(config)

	if appConfig.Recovery {
		r.App.Use(recover.New(recover.Config{
			EnableStackTrace: true,
		}))
	}

	if appConfig.RequestID {
		r.App.Use(requestid.New())
	}

	if appConfig.Logger {
		r.App.Use(logger.New(logger.Config{
			Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
			Output: log.GetWriter(),
		}))
	}

	if appConfig.Cors {
		r.App.Use(cors.New())
	}

	if appConfig.Swagger {
		r.App.Add(http.MethodGet, "/swagger/*", swagger.HandlerDefault)
	}

	if appConfig.Metrics {
		metrics := collector.New(r.appName)
		metrics.RegisterAt(r.App, "/metrics")
		r.App.Use(metrics.Middleware)
	}

	if appConfig.Static != nil {
		r.App.Static(appConfig.Static.Prefix, appConfig.Static.Root)
	}

	if appConfig.ApplicationWarmup == nil {
		r.applicationInitialized = true
	}

	if len(appConfig.Workers) > 0 {
		r.hostedServices = appConfig.Workers
	}

	r.App.Use(healthcheck.New(healthcheck.Config{
		ReadinessEndpoint: "/ping",
		ReadinessProbe: func(ctx *fiber.Ctx) bool {
			return r.isReady()
		},
	}))

	r.appConfig = appConfig
}

func (r *FastHTTPServer) Shutdown() error {
	return r.App.Shutdown()
}

func (r *FastHTTPServer) Join() error {
	for err := range r.errChan {
		return err
	}

	return nil
}

func (r *FastHTTPServer) On(application application.IApplication) {
	if application != nil {
		application.RegisterServer(r)
		r.application = application
	}
}

func (r *FastHTTPServer) AddHostedService(hostedService services.IHostedService) {
	r.hostedServices = append(r.hostedServices, hostedService)
}

func (r *FastHTTPServer) GetHostedServices() []services.IHostedService {
	return r.hostedServices
}

func (r *FastHTTPServer) RegisterRoutes(routes []routing.Route) {
	for i := 0; i < len(routes); i++ {
		route := routes[i]
		r.Add(route.Method, routes[i].Path, func(ctx *fiber.Ctx) error {
			return route.Action(&routing.HTTPContext{
				Ctx: ctx,
			})
		})
	}
}

func (r *FastHTTPServer) isReady() bool {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.applicationInitialized
}

func (r *FastHTTPServer) setIsReady(value bool) {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.applicationInitialized = value
}

func (r *FastHTTPServer) registerSubscribers(ctx context.Context) {
	r.onListenSubscribers = append(r.onListenSubscribers, subscriptions.NewApplicationWarmupListener(r))
	r.onListenSubscribers = append(r.onListenSubscribers, subscriptions.NewHostedServiceListener(r))

	for i := 0; i < len(r.onListenSubscribers); i++ {
		subscriber := r.onListenSubscribers[i]
		if subscriber.MustSubscribe() {
			r.onListen.Subscribe(ctx, subscriber)
		}
	}
}
