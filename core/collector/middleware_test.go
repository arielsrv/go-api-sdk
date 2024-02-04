package collector_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arielsrv/go-sdk-api/core/collector"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestMiddleware(t *testing.T) {
	app := fiber.New()
	metrics := collector.New("test-service")
	metrics.RegisterAt(app, "/metrics")
	app.Use(metrics.Middleware)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})
	app.Get("/error/:type", func(ctx *fiber.Ctx) error {
		switch ctx.Params("type") {
		case "fiber":
			return fiber.ErrBadRequest
		default:
			return fiber.ErrInternalServerError
		}
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusOK {
		t.Fail()
	}

	req = httptest.NewRequest(http.MethodGet, "/error/fiber", nil)
	resp, _ = app.Test(req, -1)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fail()
	}

	req = httptest.NewRequest(http.MethodGet, "/error/unknown", nil)
	resp, _ = app.Test(req, -1)
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fail()
	}

	req = httptest.NewRequest(http.MethodGet, "/metrics", nil)
	resp, _ = app.Test(req, -1)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err)
		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	got := string(body)
	want := `http_requests_total{application="test-service",method="GET",path="/",status_code="200"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `http_requests_total{application="test-service",method="GET",path="/error/:type",status_code="400"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `http_requests_total{application="test-service",method="GET",path="/error/:type",status_code="500"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `http_request_duration_seconds_count{application="test-service",method="GET",path="/",status_code="200"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `http_requests_in_progress_total{application="test-service",method="GET"} 0`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}
}

func TestMiddlewareOnRoute(t *testing.T) {
	app := fiber.New()
	metrics := collector.New("test-route")
	prefix := "/prefix/path"
	app.Route(prefix, func(route fiber.Router) {
		metrics.RegisterAt(route, "/metrics")
	}, "Prefixed Route")
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})
	app.Get("/error/:type", func(ctx *fiber.Ctx) error {
		switch ctx.Params("type") {
		case "fiber":
			return fiber.ErrBadRequest
		default:
			return fiber.ErrInternalServerError
		}
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusOK {
		t.Fail()
	}

	req = httptest.NewRequest(http.MethodGet, "/error/fiber", nil)
	resp, _ = app.Test(req, -1)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fail()
	}

	req = httptest.NewRequest(http.MethodGet, "/error/unknown", nil)
	resp, _ = app.Test(req, -1)
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fail()
	}

	req = httptest.NewRequest(http.MethodGet, prefix+"/metrics", nil)
	resp, _ = app.Test(req, -1)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err)
		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	got := string(body)
	want := `http_requests_total{application="test-service",method="GET",path="/",status_code="200"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `http_requests_total{application="test-service",method="GET",path="/error/:type",status_code="400"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `http_requests_total{application="test-service",method="GET",path="/error/:type",status_code="500"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `http_request_duration_seconds_count{application="test-service",method="GET",path="/",status_code="200"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `http_requests_in_progress_total{application="test-service",method="GET"} 0`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}
}

func TestMiddlewareWithServiceName(t *testing.T) {
	app := fiber.New()

	metrics := collector.NewWith("unique-service", "my_service_with_name", "http")
	metrics.RegisterAt(app, "/metrics")
	app.Use(metrics.Middleware)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusOK {
		t.Fail()
	}

	req = httptest.NewRequest(http.MethodGet, "/metrics", nil)
	resp, _ = app.Test(req, -1)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err)
		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	got := string(body)
	want := `my_service_with_name_http_requests_total{application="unique-service",method="GET",path="/",status_code="200"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `my_service_with_name_http_request_duration_seconds_count{application="unique-service",method="GET",path="/",status_code="200"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `my_service_with_name_http_requests_in_progress_total{application="unique-service",method="GET"} 0`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}
}

func TestMiddlewareWithLabels(t *testing.T) {
	app := fiber.New()

	constLabels := map[string]string{
		"customkey1": "customvalue1",
		"customkey2": "customvalue2",
	}
	metrics := collector.NewWithLabels(constLabels, "my_service", "http")
	metrics.RegisterAt(app, "/metrics")
	app.Use(metrics.Middleware)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusOK {
		t.Fail()
	}

	req = httptest.NewRequest(http.MethodGet, "/metrics", nil)
	resp, _ = app.Test(req, -1)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err)
		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	got := string(body)
	want := `my_service_http_requests_total{customkey1="customvalue1",customkey2="customvalue2",method="GET",path="/",status_code="200"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `my_service_http_request_duration_seconds_count{customkey1="customvalue1",customkey2="customvalue2",method="GET",path="/",status_code="200"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `my_service_http_requests_in_progress_total{customkey1="customvalue1",customkey2="customvalue2",method="GET"} 0`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}
}

func TestMiddlewareWithBasicAuth(t *testing.T) {
	app := fiber.New()

	metrics := collector.New("basic-auth")
	metrics.RegisterAt(app, "/metrics", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"prometheus": "password",
		},
	}))

	app.Use(metrics.Middleware)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusOK {
		t.Fail()
	}

	req = httptest.NewRequest(http.MethodGet, "/metrics", nil)
	resp, _ = app.Test(req, -1)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fail()
	}

	req.SetBasicAuth("prometheus", "password")
	resp, _ = app.Test(req, -1)
	if resp.StatusCode != http.StatusOK {
		t.Fail()
	}
}

func TestMiddlewareWithCustomRegistry(t *testing.T) {
	app := fiber.New()
	registry := prometheus.NewRegistry()

	srv := httptest.NewServer(promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	t.Cleanup(srv.Close)

	promfiber := collector.NewWithRegistry(registry, "unique-service", "my_service_with_name", "http", nil)
	app.Use(promfiber.Middleware)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fail()
	}
	if resp.StatusCode != http.StatusOK {
		t.Fail()
	}

	resp, err = srv.Client().Get(srv.URL)
	if err != nil {
		t.Fail()
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
	if resp.StatusCode != http.StatusOK {
		t.Fail()
	}
	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			t.Error(closeErr)
		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	got := string(body)
	want := `my_service_with_name_http_requests_total{application="unique-service",method="GET",path="/",status_code="200"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `my_service_with_name_http_request_duration_seconds_count{application="unique-service",method="GET",path="/",status_code="200"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `my_service_with_name_http_requests_in_progress_total{application="unique-service",method="GET"} 0`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}
}
