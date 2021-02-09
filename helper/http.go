package helper

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/patcharp/golib/crypto"
	"github.com/patcharp/golib/server"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"runtime"
)

/**
 * General response
 */
func HttpResponse(c echo.Context, code int, result server.Result) error {
	return c.JSON(code, result)
}

// 200
func HttpOk(c echo.Context, data interface{}) error {
	count := 0
	if data != nil {
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			count = reflect.ValueOf(data).Len()
		}
	}
	return c.JSON(http.StatusOK, server.Result{
		Message: "success",
		Data:    data,
		Count:   count,
	})
}

// 200
func HttpOkWithTotal(c echo.Context, data interface{}, total int) error {
	count := 0
	if data != nil {
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			count = reflect.ValueOf(data).Len()
		}
	}
	return c.JSON(http.StatusOK, server.Result{
		Message: "success",
		Data:    data,
		Total:   total,
		Count:   count,
	})
}

// 201
func HttpCreated(c echo.Context, data interface{}) error {
	return HttpResponse(c, http.StatusCreated, server.Result{
		Message: "success",
		Data:    data,
	})
}

// 204
func HttpNoContent(c echo.Context) error {
	return HttpResponse(c, http.StatusNoContent, server.Result{
		Message: "success",
	})
}

/**
 * Client request
 */
func HttpInvalidRequest(c echo.Context, code int, err error, msg interface{}) error {
	// TODO: Generate error code from error and return error code to frontend
	errCode := crypto.GenSecretString(8)
	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	logrus.Errorln(fmt.Sprintf(
		"[CLIENT%d] %s (%s):%d error #%s -> %v",
		code,
		funcName,
		file,
		line,
		errCode,
		err,
	))
	return HttpResponse(c, code, server.Result{
		Message: msg,
		Error:   fmt.Sprintf("#%s: %v", errCode, err),
	})
}

// 400
func HttpErrBadRequest(c echo.Context) error {
	return HttpInvalidRequest(c, http.StatusBadRequest, nil, "invalid request")
}

// 401
func HttpErrUnAuthorize(c echo.Context, err error) error {
	return HttpInvalidRequest(c, http.StatusUnauthorized, err, "unauthorized")
}

// 403
func HttpErrForbidden(c echo.Context) error {
	return HttpInvalidRequest(c, http.StatusForbidden, nil, "forbidden")
}

// 404
func HttpErrNotFound(c echo.Context) error {
	return HttpInvalidRequest(c, http.StatusNotFound, nil, "request not found")
}

// 409
func HttpErrConflict(c echo.Context) error {
	return HttpInvalidRequest(c, http.StatusConflict, nil, "requested was conflict")
}

/**
 * Server error
 */
func HttpServerError(c echo.Context, code int, err error, msg interface{}) error {
	// TODO: Generate error code from error and return error code to frontend
	errCode := crypto.GenSecretString(8)
	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	logrus.Errorln(fmt.Sprintf(
		"[SERVER%d] %s (%s):%d error #%s -> %v",
		code,
		funcName,
		file,
		line,
		errCode,
		err,
	))
	return HttpResponse(c, code, server.Result{
		Message: msg,
		Error:   fmt.Sprintf("#%s: %v", errCode, err),
	})
}

// 500
func HttpErrServerError(c echo.Context, err error, msg interface{}) error {
	return HttpServerError(c, http.StatusInternalServerError, err, "server error")
}

// 503
func HttpServiceUnavailableError(c echo.Context, err error, msg interface{}) error {
	return HttpServerError(c, http.StatusServiceUnavailable, err, "destination service unavailable")
}
