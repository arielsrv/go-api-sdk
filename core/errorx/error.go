package errorx

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/oops"
	log "gitlab.com/iskaypetcom/digital/sre/tools/dev/go-logger"
)

// swagger:response Error
type Error struct {
	Type     string `example:"https://example.com/probs/out-of-credit"        json:"type,omitempty"`
	Title    string `example:"You do not have enough credit."                 json:"title,omitempty"`
	Status   int    `example:"403"                                            json:"status,omitempty"`
	Detail   string `example:"Your current balance is 30, but that costs 50." json:"detail,omitempty"`
	Instance string `example:"/account/12345/msgs/abc"                        json:"instance,omitempty"`

	OopsError error `json:"-"`
}

// ThrowErr feel free to use in place (service, handlers, controllers, etc.).
func ThrowErr(statusCode int, err error) *Error {
	return &Error{
		Status:    statusCode,
		Title:     http.StatusText(statusCode),
		OopsError: oops.Wrap(err),
	}
}

func (e Error) Error() string {
	bytes, err := json.Marshal(e)
	if err != nil {
		return e.Title
	}

	return string(bytes)
}

// OnErrorHandler is the default error handler.
func OnErrorHandler(ctx *fiber.Ctx, err error) error {
	var e = new(Error)
	var fiberError *fiber.Error
	var apiError *Error

	switch {
	case errors.As(err, &fiberError):
		{
			e.Status = fiberError.Code
			e.Title = http.StatusText(e.Status)
			wErr := oops.Wrap(err)
			var oopsErr oops.OopsError
			if errors.As(wErr, &oopsErr) {
				e.Detail = oopsErr.Stacktrace()
			}
		}
	case errors.As(err, &apiError):
		{
			e.Status = apiError.Status
			e.Title = apiError.Title
			e.Type = apiError.Type
			e.Instance = apiError.Instance
			e.Detail = apiError.Detail
			var oopsErr oops.OopsError
			if apiError.OopsError != nil {
				if errors.As(apiError.OopsError, &oopsErr) {
					e.Detail = oopsErr.Stacktrace()
				}
			}
		}
	default:
		e.Status = http.StatusInternalServerError
		e.Title = http.StatusText(e.Status)
		oopsErr := oops.Wrap(err)
		var o oops.OopsError
		if errors.As(oopsErr, &o) {
			e.Detail = o.Stacktrace()
		}
	}

	if e.Status >= http.StatusInternalServerError {
		log.Error(e.Detail)
	}

	e.Instance = ctx.Path()
	ctx.Status(e.Status)
	ctx.Set(fiber.HeaderContentType, "application/problem+json; charset=utf-8")

	bytes, mErr := json.Marshal(e)
	if mErr != nil {
		return ctx.JSON(e)
	}

	return ctx.SendString(string(bytes))
}
