package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"os"
	"os/signal"
)

type Config struct {
	Host   string
	Port   string
	Prefix string
	Prod   bool
	// Fiber Config
	Config *fiber.Config
}

type Server struct {
	config Config
	app    *fiber.App
}

func New(config Config) Server {
	s := Server{
		config: config,
	}
	if config.Config != nil {
		s.app = fiber.New(*config.Config)
	} else {
		s.app = fiber.New()
	}
	s.app.Use(recover.New())
	return s
}

func (s *Server) Run() error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		_ = s.app.Shutdown()
	}()
	return s.app.Listen(fmt.Sprintf("%s:%s", s.config.Host, s.config.Port))
}

func (s *Server) App() *fiber.App {
	return s.app
}

// DefaultErrorHandler that process return errors from handlers
var DefaultErrorHandler = func(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(code).JSON(fiber.Map{"error": err.Error()})
}
