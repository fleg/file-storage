package router

import (
	"net/http"

	"file-storage/internal/services"

	"github.com/labstack/echo/v4"
)

type (
	HealthController struct {
		healthService *services.HealthService
	}

	CheckResponse struct {
		Status string `json:"status"`
		Time   int64  `json:"time"`
	}
)

func (hh *HealthController) Health(c echo.Context) error {
	health := hh.healthService.GetHealth(c.Request().Context())

	return c.JSON(http.StatusOK, CheckResponse{
		Status: health.Status,
		Time:   health.Time,
	})
}

func (hh *HealthController) Ready(c echo.Context) error {
	ready, err := hh.healthService.GetReadiness(c.Request().Context())

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, CheckResponse{
		Status: ready.Status,
		Time:   ready.Time,
	})
}

func NewHealthController(hs *services.HealthService) *HealthController {
	return &HealthController{healthService: hs}
}
