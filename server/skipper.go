package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type SkipperPath struct {
	Prefix string
	Paths  map[string]bool
}

func (s *SkipperPath) Add(path string, method string) {
	if s.Prefix != "" {
		path = fmt.Sprintf("%s%s", s.Prefix, path)
	}
	s.Paths[s.key(path, method)] = true
}

func (s *SkipperPath) Delete(path string, method string) {
	if s.Prefix != "" {
		path = fmt.Sprintf("%s%s", s.Prefix, path)
	}
	delete(s.Paths, s.key(path, method))
}

func (s *SkipperPath) key(path string, method string) string {
	return fmt.Sprintf("%s|%s", method, path)
}

func (s *SkipperPath) Test(ctx *fiber.Ctx) bool {
	if active, ok := s.Paths[s.key(ctx.Path(), ctx.Method())]; ok && active {
		return true
	}
	return false
}

func NewSkipperPath(prefix string) SkipperPath {
	return SkipperPath{
		Prefix: prefix,
		Paths:  map[string]bool{},
	}
}
