package helper

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/patcharp/golib/v2/server"
	"github.com/patcharp/golib/v2/util"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"runtime"
)

/**
 * General response
 */
func HttpResponse(ctx *fiber.Ctx, code int, result server.Result) error {
	return ctx.Status(code).JSON(result)
}

// 200
func HttpOk(ctx *fiber.Ctx, data interface{}) error {
	count := 0
	if data != nil {
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			count = reflect.ValueOf(data).Len()
		}
	}
	return ctx.JSON(server.Result{
		Message: "success",
		Data:    data,
		Count:   count,
	})
}

// 200
func HttpOkWithTotal(ctx *fiber.Ctx, data interface{}, total int) error {
	count := 0
	if data != nil {
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			count = reflect.ValueOf(data).Len()
		}
	}
	return ctx.JSON(server.Result{
		Message: "success",
		Data:    data,
		Total:   total,
		Count:   count,
	})
}

// 201
func HttpCreated(ctx *fiber.Ctx, data interface{}) error {
	return HttpResponse(ctx, http.StatusCreated, server.Result{
		Message: "success",
		Data:    data,
	})
}

// 204
func HttpNoContent(ctx *fiber.Ctx) error {
	return HttpResponse(ctx, http.StatusNoContent, server.Result{
		Message: "success",
	})
}

/**
 * Blob
 */
func HttpBlob(ctx *fiber.Ctx, blob []byte) error {
	return HttpBlobWithCode(ctx, blob, http.StatusOK)
}

func HttpBlobWithCode(ctx *fiber.Ctx, blob []byte, code int) error {
	// http.DetectContentType(blob)
	return ctx.Status(code).Send(blob)
}

/**
 * Redirect
 */
func HttpRedirect(ctx *fiber.Ctx, url string) error {
	return HttpRedirectWithCode(ctx, url, http.StatusFound)
}

func HttpRedirectWithCode(ctx *fiber.Ctx, url string, code int) error {
	return ctx.Status(code).Redirect(url)
}

/**
 * Client request
 */
func HttpInvalidRequest(ctx *fiber.Ctx, code int, err error, msg interface{}) error {
	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	if util.GetEnv("HTTP_DEBUG", "false") == "true" {
		logrus.Errorln(fmt.Sprintf(
			"[CLIENT%d] %s (%s):%d error -> %v",
			code,
			funcName,
			file,
			line,
			err,
		))
	}
	return HttpResponse(ctx, code, server.Result{
		Error: msg,
	})
}

// 400
func HttpErrBadRequest(ctx *fiber.Ctx) error {
	return HttpInvalidRequest(ctx, http.StatusBadRequest, nil, "invalid request")
}

// 401
func HttpErrUnAuthorize(ctx *fiber.Ctx, err error) error {
	return HttpInvalidRequest(ctx, http.StatusUnauthorized, err, "unauthorized")
}

// 403
func HttpErrForbidden(ctx *fiber.Ctx) error {
	return HttpInvalidRequest(ctx, http.StatusForbidden, nil, "forbidden")
}

// 404
func HttpErrNotFound(ctx *fiber.Ctx) error {
	return HttpInvalidRequest(ctx, http.StatusNotFound, nil, "request not found")
}

// 409
func HttpErrConflict(ctx *fiber.Ctx) error {
	return HttpInvalidRequest(ctx, http.StatusConflict, nil, "requested was conflict")
}

/**
 * Server error
 */
func HttpServerError(ctx *fiber.Ctx, code int, err error, msg interface{}) error {
	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	logrus.Errorln(fmt.Sprintf(
		"[SERVER%d] %s (%s):%d error -> %v",
		code,
		funcName,
		file,
		line,
		err,
	))
	return HttpResponse(ctx, code, server.Result{
		Error: msg,
	})
}

// 500
func HttpErrServerError(ctx *fiber.Ctx, err error, msg interface{}) error {
	return HttpServerError(ctx, http.StatusInternalServerError, err, "server error")
}

// 503
func HttpServiceUnavailableError(ctx *fiber.Ctx, err error, msg interface{}) error {
	return HttpServerError(ctx, http.StatusServiceUnavailable, err, "destination service unavailable")
}
