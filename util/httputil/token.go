package httputil

import (
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"strings"
)

// Extract token from http header
func GetTokenFromHeader(c echo.Context, tokenType string, header string) string {
	token := c.Request().Header.Get(header)
	token = strings.TrimSpace(token)
	tokenType += " "
	if token == "" || len(token) < (len(tokenType)+1) || strings.ToLower(token[:len(tokenType)]) != strings.ToLower(tokenType) {
		return ""
	}
	token = strings.TrimSpace(token[len(tokenType):])
	return token
}

func GetQueryId(c echo.Context, qName string, id *uuid.UUID) error {
	var err error
	if *id, err = uuid.FromString(c.Param(qName)); err != nil {
		return err
	}
	return nil
}
