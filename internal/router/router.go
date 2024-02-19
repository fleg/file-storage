package router

import (
	_ "file-storage/docs"
	"file-storage/internal/config"
	"file-storage/internal/server"
	"file-storage/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
)

//	@title			file-storage service
//	@version		1.0
//	@description	Simple file storage service
//	@contact.name	Aleksey Timchenko
//	@contact.url	https://flegdev.ru
//	@contact.email	aleksey.timchenqo@gmail.com
//	@BasePath		/
func Bind(srv *server.Server, c *config.Config, hs *services.HealthService, fs *services.FilesService) {
	bindHealth(srv.Echo, hs)
	bindFiles(srv.Echo, fs)

	if c.Swagger.IsEnabled {
		bindSwagger(srv.Echo)
	}
}

func bindHealth(e *echo.Echo, hs *services.HealthService) {
	hc := NewHealthController(hs)

	e.GET("/health", hc.Health)
	e.GET("/ready", hc.Ready)
}

func bindFiles(e *echo.Echo, fs *services.FilesService) {
	fc := NewFilesController(fs)

	e.POST("/upload", fc.Upload)
	e.GET("/download/:id", fc.Download)
	e.GET("/files/:id", fc.GetOne)
	e.DELETE("/files/:id", fc.Unlink)
}

func bindSwagger(e *echo.Echo) {
	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.GET("/docs", func(c echo.Context) error {
		return c.Redirect(301, "/docs/index.html")
	})
}
