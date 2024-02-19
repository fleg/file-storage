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
		Status string `json:"status"  example:"OK"`
		Time   int64  `json:"time" example:"1708353333135"`
	}
)

// Health godoc
// @Summary      Health route
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  router.CheckResponse
// @Router       /health [get]
func (hc *HealthController) Health(c echo.Context) error {
	health := hc.healthService.GetHealth(c.Request().Context())

	return c.JSON(http.StatusOK, CheckResponse{
		Status: health.Status,
		Time:   health.Time,
	})
}

// Ready godoc
// @Summary      Ready route
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  router.CheckResponse
// @Router       /ready [get]
func (hc *HealthController) Ready(c echo.Context) error {
	ready, err := hc.healthService.GetReadiness(c.Request().Context())

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
