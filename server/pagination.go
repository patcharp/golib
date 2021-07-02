package server

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

const (
	DefaultQueryPage = 1
	DefaultQuerySize = 20
	MaxQuerySize     = 100
)

type Pagination struct {
	Page int
	Size int
}

func GetPagination(ctx *fiber.Ctx) Pagination {
	var p Pagination
	p.Page = atoi(ctx.Query("page"), DefaultQueryPage)
	p.Size = atoi(ctx.Query("size"), DefaultQuerySize)
	if p.Size > MaxQuerySize {
		p.Size = MaxQuerySize
	}
	return p
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.Size
}

func atoi(s string, v int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return v
	}
	return i
}
