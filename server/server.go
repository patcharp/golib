package server

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/patcharp/golib/v2/util"
	"github.com/patcharp/golib/v2/util/httputil"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"time"
)

type Config struct {
	Host        string
	Port        string
	Config      *fiber.Config
	HealthCheck bool
	RequestId   bool
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

	s.app.Use(logger.New(logger.Config{
		Format:     "${locals:requestid} - ${ip} - ${method} ${path} ${status} - ${latency}\n",
		TimeZone:   "Asia/Bangkok",
		TimeFormat: time.ANSIC,
		Next: func(c *fiber.Ctx) bool {
			// no log for health check
			if c.Path() == "/api/-/health" {
				return true
			}
			return false
		},
	}))

	s.app.Use(recover.New(recover.Config{
		EnableStackTrace:  true,
		StackTraceHandler: fiberStackTraceHandler,
	}))

	if config.HealthCheck {
		s.EnableHealthCheck()
	}

	if config.RequestId {
		s.EnableRequestId()
	}

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

func (s *Server) EnableRequestId() {
	s.config.RequestId = true
	s.app.Use(requestid.New(requestid.Config{
		Generator: func() string {
			return ksuid.New().String()
		},
	}))
}

func (s *Server) EnableHealthCheck() {
	s.config.HealthCheck = true
	s.app.Get("/api/-/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).SendString("ok")
	})
}

func (s *Server) App() *fiber.App {
	return s.app
}

// DefaultErrorHandler that process return errors from handlers
var DefaultErrorHandler = func(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}
	if code >= fiber.StatusInternalServerError {
		logrus.Error("[PANIC] "+fmt.Sprintf(" [%s] ", ctx.IP())+ctx.Route().Method+ctx.Route().Path+" -> ", err)
	}
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return ctx.Status(code).JSON(fiber.Map{"error": err.Error()})
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

func fiberStackTraceHandler(_ *fiber.Ctx, e interface{}) {
	logrus.Error(fmt.Sprintf("[PANIC] %v\n%s\n", e, debug.Stack()))
}
