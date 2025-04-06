package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

// HealthCheck is a special API that will be called periodically by docker for their health checking functionality.
func (integrator APIIntegrator) HealthCheck(c echo.Context) error {
	integrator.Logger.NewInfo(fmt.Sprintf("[Integrator][HealthCheck] Health Check called at time %s", time.Now().String()))
	return c.JSON(http.StatusOK, nil)
}
