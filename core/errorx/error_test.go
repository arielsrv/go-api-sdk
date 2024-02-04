package errorx_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/errorx"
)

func TestNewError(t *testing.T) {
	actual := errorx.ThrowErr(http.StatusInternalServerError, errors.New("nil reference"))
	assert.NotNil(t, actual)
	assert.Equal(t, http.StatusInternalServerError, actual.Status)
	require.ErrorAs(t, actual.OopsError, &oops.OopsError{})
}

func TestErrorHandler(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)
	err := errorx.OnErrorHandler(ctx, errors.New("api server error"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, ctx.Context().Response.StatusCode())

	var apiError errorx.Error
	err = json.Unmarshal(ctx.Response().Body(), &apiError)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, apiError.Status)
	assert.Equal(t, http.StatusText(http.StatusInternalServerError), apiError.Title)
}

func TestErrorHandler_FiberError(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)
	err := errorx.OnErrorHandler(ctx, fiber.NewError(http.StatusInternalServerError, "api server error"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, ctx.Context().Response.StatusCode())

	var apiError errorx.Error
	err = json.Unmarshal(ctx.Response().Body(), &apiError)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, apiError.Status)
	assert.Equal(t, http.StatusText(http.StatusInternalServerError), apiError.Title)
}

func TestErrorHandler_ApiError(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)
	err := errorx.OnErrorHandler(ctx, errorx.ThrowErr(http.StatusInternalServerError, errors.New("api server error")))
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, ctx.Context().Response.StatusCode())

	var apiError errorx.Error
	err = json.Unmarshal(ctx.Response().Body(), &apiError)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, apiError.Status)
	assert.Equal(t, http.StatusText(http.StatusInternalServerError), apiError.Title)
}

func TestError_Error(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)
	err := errorx.OnErrorHandler(ctx, errors.New("api server error"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, ctx.Context().Response.StatusCode())

	var apiError errorx.Error
	err = json.Unmarshal(ctx.Response().Body(), &apiError)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, apiError.Status)
	assert.Contains(t, apiError.Error(), "api server error")
}
