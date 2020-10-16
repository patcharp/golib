package server

import (
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	echoPrometheus "github.com/globocom/echo-prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type Config struct {
	Host   string
	Port   string
	Prefix string
	Prod   bool
}

type Server struct {
	config Config
	ctx    *echo.Echo
}

type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Message interface{} `json:"message,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Total   int         `json:"total,omitempty"`
	Count   int         `json:"count,omitempty"`
}

func New(config Config) Server {
	s := Server{
		config: config,
		ctx:    echo.New(),
	}
	s.ctx.Use(middleware.Recover())
	s.ctx.Use(middleware.RemoveTrailingSlash())
	// Cr. https://echo.labstack.com/middleware/request-id
	s.ctx.Use(middleware.RequestID())
	// Cr. https://echo.labstack.com/middleware/secure
	s.ctx.Use(middleware.Secure())
	s.ctx.HTTPErrorHandler = s.serverErrorHandler
	listenAddr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	s.ctx.Server = &http.Server{
		Addr:    listenAddr,
		Handler: s.ctx,
	}
	if !s.config.Prod {
		s.ctx.Use(s.logger())
	}
	return s
}

func (s *Server) Run() error {
	if err := gracehttp.Serve(s.ctx.Server); err != nil {
		return err
	}
	return nil
}

func (s *Server) Ctx() *echo.Echo {
	return s.ctx
}

func (s *Server) EnableMetrics(metricsPath string, nameSpace string) error {
	s.ctx.Use(echoPrometheus.MetricsMiddlewareWithConfig(echoPrometheus.Config{
		Namespace: nameSpace,
		Buckets: []float64{
			0.0005, // 0.5ms
			0.001,  // 1ms
			0.005,  // 5ms
			0.01,   // 10ms
			0.05,   // 50ms
			0.1,    // 100ms
			0.5,    // 500ms
			1,      // 1s
			2,      // 2s
			5,      // 5s
			10,     // 10s
		},
	}))
	s.ctx.GET(metricsPath, echo.WrapHandler(promhttp.Handler()))
	return nil
}

func (s *Server) EnableCORS(allowOrigins []string, allowHeaders []string) error {
	if len(allowOrigins) == 0 {
		allowOrigins = []string{"*"}
	}
	if len(allowHeaders) == 0 {
		allowHeaders = []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization}
	}
	s.ctx.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowHeaders: allowOrigins,
	}))
	return nil
}

func (s *Server) serverErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
		if he.Internal != nil {
			err = fmt.Errorf("%v, %v", err, he.Internal)
		}
	} else if s.ctx.Debug {
		msg = err.Error()
	} else {
		msg = http.StatusText(code)
	}
	if _, ok := msg.(string); ok {
		msg = map[string]interface{}{"error": msg}
	}
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead { // Issue #608
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, msg)
		}
		if err != nil {
			s.ctx.Logger.Error(err)
		}
	}
}

func (s *Server) logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			var err error
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}
			reqSize := req.Header.Get(echo.HeaderContentLength)
			if reqSize == "" {
				reqSize = "0"
			}
			log.Infof("%s %s [%v] %s %-7s %s %3d %s %s %13v %s %s",
				id,
				c.RealIP(),
				stop.Format(time.RFC3339),
				req.Host,
				req.Method,
				req.RequestURI,
				res.Status,
				reqSize,
				strconv.FormatInt(res.Size, 10),
				stop.Sub(start).String(),
				req.Referer(),
				req.UserAgent(),
			)
			return err
		}
	}
}
