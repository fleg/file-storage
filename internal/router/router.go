package router

import (
	"file-storage/internal/config"
	"file-storage/internal/server"
	"file-storage/internal/services"

	"github.com/labstack/echo/v4"
)

func Bind(srv *server.Server, c *config.Config, hs *services.HealthService, fs *services.FilesService) {
	bindHealth(srv.Echo, hs)
	bindFiles(srv.Echo, fs)
}

func bindHealth(e *echo.Echo, hs *services.HealthService) {
	hc := NewHealthController(hs)

	e.GET("/health", hc.Health)
	e.GET("/ready", hc.Ready)
}

func bindFiles(e *echo.Echo, fs *services.FilesService) {
	fc := NewFilesController(fs)

	e.POST("/upload", fc.Upload)
}
