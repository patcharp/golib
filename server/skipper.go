package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

type SkipperPath struct {
	Prefix string
	Paths  map[string]bool
}

func (s *SkipperPath) Add(path string, method string) {
	s.Paths[s.key(path, method)] = true
}

func (s *SkipperPath) Delete(path string, method string) {
	delete(s.Paths, s.key(path, method))
}

func (s *SkipperPath) key(path string, method string) string {
	return fmt.Sprintf("%s%s", method, path)
}

func (s *SkipperPath) Test(c echo.Context) bool {
	if active, ok := s.Paths[s.key(c.Path(), c.Request().Method)]; ok && active {
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
