package server

import (
	"github.com/labstack/echo/v4"
	"github.com/patcharp/golib/util"
	"net/http"
	"strings"
)

func EnableCORS(c echo.Context, allowOrigins []string, allowHeaders []string) error {
	defaultAllowHeaders := []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization}
	if len(allowOrigins) == 0 {
		allowOrigins = []string{"*"}
	}
	for _, h := range defaultAllowHeaders {
		if !util.Contains(allowHeaders, h) {
			allowHeaders = append(allowHeaders, h)
		}
	}
	c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, strings.Join(allowOrigins, ","))
	c.Response().Header().Set(echo.HeaderAccessControlAllowHeaders, strings.Join(allowHeaders, ","))
	return c.NoContent(http.StatusNoContent)
}
