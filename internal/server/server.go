package server

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/brpaz/echozap"

	"file-storage/internal/logger"
)

type (
	Server struct {
		Echo *echo.Echo
	}

	StartOptions struct {
		Port uint
	}
)

func New(l *logger.Logger) *Server {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = NewHTTPErrorHandler()

	e.Use(echozap.ZapLogger(l.Desugar()))
	e.Use(middleware.RequestID())

	return &Server{
		Echo: e,
	}
}

func (s *Server) Start(options *StartOptions) error {
	return s.Echo.Start(fmt.Sprintf("0.0.0.0:%d", options.Port))
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Echo.Shutdown(ctx)
}
