package httputil

import (
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/ksuid"
	"strings"
)

func GetTokenFromHeader(ctx *fiber.Ctx, tokenType string) string {
	token := ctx.Get(HeaderAuthorization)
	token = strings.TrimSpace(token)
	tokenType += " "
	if token == "" || len(token) < (len(tokenType)+1) || strings.ToLower(token[:len(tokenType)]) != strings.ToLower(tokenType) {
		return ""
	}
	token = strings.TrimSpace(token[len(tokenType):])
	return token
}

func GetTokenFromCookie(ctx *fiber.Ctx, name string) string {
	token := ctx.Cookies(name)
	return token
}

func GetApiKey(ctx *fiber.Ctx) string {
	return ctx.Get(HeaderXApiKey)
}

func GetCtxParamId(ctx *fiber.Ctx, qName string, id *uuid.UUID) error {
	var err error
	if id == nil {
		id = &uuid.UUID{}
	}
	if *id, err = uuid.FromString(ctx.Params(qName)); err != nil {
		return err
	}
	return nil
}

func GetCtxParamKId(ctx *fiber.Ctx, param string, id *ksuid.KSUID) error {
	var err error
	if id == nil {
		id = &ksuid.KSUID{}
	}
	if *id, err = ksuid.Parse(ctx.Params(param)); err != nil {
		return err
	}
	return nil
}
