package router

import (
	"file-storage/internal/config"
	"file-storage/internal/server"
	"file-storage/internal/services"

	"github.com/labstack/echo/v4"
)

func Bind(srv *server.Server, c *config.Config, hs *services.HealthService) {
	bindHealth(srv.Echo, hs)
}

func bindHealth(e *echo.Echo, hs *services.HealthService) {
	hh := NewHealthController(hs)

	e.GET("/health", hh.Health)
	e.GET("/ready", hh.Ready)
}
