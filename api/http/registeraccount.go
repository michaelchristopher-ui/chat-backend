package http

import (
	"encoding/json"
	"net/http"
	"websocket_client/api/http/structs"
	"websocket_client/internal/pkg/core/adapter/accountadapter"

	"github.com/labstack/echo"
)

// GetMessages is the handler method for the /get_messages api endpoint
func (integrator APIIntegrator) RegisterAccount(c echo.Context) error {
	req := RegisterAccountReq{}
	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, structs.ErrorRet{
			Error: err.Error(),
		})
	}
	err = integrator.AccountService.Register(accountadapter.RegisterReq{
		UserID:   req.UserID,
		Password: req.Password,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, structs.ErrorRet{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, nil)
}

type RegisterAccountReq struct {
	UserID   string `json:"userid"`
	Password string `json:"password"`
}
