package collector

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/errorx"
)

// FiberPrometheus ...
type FiberPrometheus struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestInFlight *prometheus.GaugeVec
	defaultURL      string
}

func create(registry prometheus.Registerer, serviceName, namespace, subsystem string, labels map[string]string) *FiberPrometheus {
	constLabels := make(prometheus.Labels)
	if serviceName != "" {
		constLabels["application"] = serviceName
	}
	for label, value := range labels {
		constLabels[label] = value
	}

	counter := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Name:        prometheus.BuildFQName(namespace, subsystem, "requests_total"),
			Help:        "Count all http requests by status code, method and path.",
			ConstLabels: constLabels,
		},
		[]string{"status_code", "method", "path"},
	)
	histogram := promauto.With(registry).NewHistogramVec(prometheus.HistogramOpts{
		Name:        prometheus.BuildFQName(namespace, subsystem, "request_duration_seconds"),
		Help:        "Duration of all HTTP requests by status code, method and path.",
		ConstLabels: constLabels,
		Buckets: []float64{
			0.000000001, // 1ns
			0.000000002,
			0.000000005,
			0.00000001, // 10ns
			0.00000002,
			0.00000005,
			0.0000001, // 100ns
			0.0000002,
			0.0000005,
			0.000001, // 1µs
			0.000002,
			0.000005,
			0.00001, // 10µs
			0.00002,
			0.00005,
			0.0001, // 100µs
			0.0002,
			0.0005,
			0.001, // 1ms
			0.002,
			0.005,
			0.01, // 10ms
			0.02,
			0.05,
			0.1, // 100 ms
			0.2,
			0.5,
			1.0, // 1s
			2.0,
			5.0,
			10.0, // 10s
			15.0,
			20.0,
			30.0,
		},
	},
		[]string{"status_code", "method", "path"},
	)

	gauge := promauto.With(registry).NewGaugeVec(prometheus.GaugeOpts{
		Name:        prometheus.BuildFQName(namespace, subsystem, "requests_in_progress_total"),
		Help:        "All the requests in progress",
		ConstLabels: constLabels,
	}, []string{"method"})

	return &FiberPrometheus{
		requestsTotal:   counter,
		requestDuration: histogram,
		requestInFlight: gauge,
		defaultURL:      "/metrics",
	}
}

var collector sync.Once
var fiberPrometheus *FiberPrometheus

// New creates a new instance of FiberPrometheus middleware
// serviceName is available as a const label.
func New(serviceName string) *FiberPrometheus {
	collector.Do(func() {
		fiberPrometheus = create(prometheus.DefaultRegisterer, serviceName, "http", "", nil)
	})

	return fiberPrometheus
}

// NewWith creates a new instance of FiberPrometheus middleware but with an ability
// to pass namespace and a custom subsystem
// Here serviceName is created as a constant-label for the metrics
// Namespace, subsystem get prefixed to the metrics.
//
// For e.g. namespace = "my_app", subsystem = "http" then metrics would be
// `my_app_http_requests_total{...,application= "serviceName"}`.
func NewWith(serviceName, namespace, subsystem string) *FiberPrometheus {
	return create(prometheus.DefaultRegisterer, serviceName, namespace, subsystem, nil)
}

// NewWithLabels creates a new instance of FiberPrometheus middleware but with an ability
// to pass namespace and a custom subsystem
// Here labels are created as a constant-labels for the metrics
// Namespace, subsystem get prefixed to the metrics.
//
// For e.g. namespace = "my_app", subsystem = "http" and labels = map[string]string{"key1": "value1", "key2":"value2"}
// then then metrics would become
// `my_app_http_requests_total{...,key1= "value1", key2= "value2" }`.
func NewWithLabels(labels map[string]string, namespace, subsystem string) *FiberPrometheus {
	return create(prometheus.DefaultRegisterer, "", namespace, subsystem, labels)
}

// NewWithRegistry creates a new instance of FiberPrometheus middleware but with an ability
// to pass a custom registry, serviceName, namespace, subsystem and labels
// Here labels are created as a constant-labels for the metrics
// Namespace, subsystem get prefixed to the metrics.
//
// For e.g. namespace = "my_app", subsystem = "http" and labels = map[string]string{"key1": "value1", "key2":"value2"}
// then then metrics would become
// `my_app_http_requests_total{...,key1= "value1", key2= "value2" }`.
func NewWithRegistry(registry prometheus.Registerer, serviceName, namespace, subsystem string, labels map[string]string) *FiberPrometheus {
	return create(registry, serviceName, namespace, subsystem, labels)
}

// RegisterAt will register the prometheus handler at a given URL.
func (ps *FiberPrometheus) RegisterAt(app fiber.Router, url string, handlers ...fiber.Handler) {
	ps.defaultURL = url

	var h = append(handlers, adaptor.HTTPHandler(promhttp.Handler()))
	app.Get(ps.defaultURL, h...)
}

// Middleware is the actual default middleware implementation.
func (ps *FiberPrometheus) Middleware(ctx *fiber.Ctx) error {
	start := time.Now()
	method := ctx.Route().Method

	if ctx.Route().Path == ps.defaultURL {
		return ctx.Next()
	}

	ps.requestInFlight.WithLabelValues(method).Inc()
	defer func() {
		ps.requestInFlight.WithLabelValues(method).Dec()
	}()

	err := ctx.Next()
	// initialize with default error code
	// https://docs.gofiber.io/guide/error-handling
	status := fiber.StatusInternalServerError
	if err != nil {
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			// Get correct error code from fiber.Error type
			status = fiberErr.Code
		}
		var coreErr *errorx.Error
		if errors.As(err, &coreErr) {
			// Get correct error code from fiber.Error type
			status = coreErr.Status
		}
	} else {
		status = ctx.Response().StatusCode()
	}

	path := ctx.Route().Path

	statusCode := strconv.Itoa(status)
	ps.requestsTotal.WithLabelValues(statusCode, method, path).Inc()

	elapsed := float64(time.Since(start).Nanoseconds()) / 1e9
	ps.requestDuration.WithLabelValues(statusCode, method, path).Observe(elapsed)

	return err
}
