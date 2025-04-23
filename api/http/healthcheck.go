package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// HealthCheck is a special API that will be called periodically by docker for their health checking functionality.
func (integrator APIIntegrator) HealthCheck(c echo.Context) error {
	integrator.Logger.NewInfo("Health Check called at time %s", time.Now().String())
	return c.JSON(http.StatusOK, nil)
}
