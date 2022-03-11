package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/patcharp/golib/v2/util"
	"github.com/patcharp/golib/v2/util/httputil"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"
)

type Config struct {
	Host string
	Port string
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

// DefaultServerAccessLog that use for enable fiber log to specific router
var DefaultServerAccessLog = func(lvl logrus.Level) fiber.Handler {
	logOut := logrus.New()
	logOut.SetLevel(lvl)
	return logger.New(logger.Config{
		Format:     "${time} ${method} ${path} - ${ip} - ${status} - ${latency}\n",
		TimeZone:   "Asia/Bangkok",
		TimeFormat: time.ANSIC,
		Output:     logOut.WriterLevel(lvl),
	})
}

// DefaultFiberConfig that use for standard fiber configuration
var DefaultFiberConfig = fiber.Config{
	Prefork:               util.GetEnv("HTTP_PRE_FORK", "false") == "true",
	ServerHeader:          util.GetEnv("HTTP_SERVER_HEADER", "GoFiber"),
	ProxyHeader:           util.GetEnv("HTTP_PROXY_HEADER", httputil.HeaderXForwardedFor),
	ReduceMemoryUsage:     util.GetEnv("HTTP_REDUCE_MEMORY_USAGE", "true") == "true",
	DisableStartupMessage: util.GetEnv("HTTP_DISABLE_STARTUP_MESSAGE", "true") == "true",
}

// ApplyStaticRoute that use for serve angular built
func ApplyStaticRoute(s *fiber.App, publicDir string) {
	staticCfg := fiber.Static{
		CacheDuration: time.Duration(0),
	}
	s.Static("/public", publicDir, staticCfg)
	s.Static("/assets", fmt.Sprintf("%s/assets", publicDir))
	s.Get("/*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile(fmt.Sprintf("%s/index.html", publicDir))
	})
}
