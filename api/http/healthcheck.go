package http

import (
	"net/http"

	"github.com/labstack/echo"
)

func (integrator APIIntegrator) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
