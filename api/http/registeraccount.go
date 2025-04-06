package http

import (
	"encoding/json"
	"fmt"
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
		integrator.Logger.NewError("[Integrator][RegisterAccount] Error when decoding, err: %s")
		return c.JSON(http.StatusInternalServerError, structs.ErrorRet{
			Error: err.Error(),
		})
	}
	//Validate password length. Since we're using BCrypt it must be 72 characters or less, starting from the 0 index.
	if len(req.Password) > 72 {
		integrator.Logger.NewInfo(fmt.Sprintf("[Integrator][RegisterAccount] Password is not shorter than 72 bytes, ID: %s", req.UserID))
		return c.JSON(http.StatusBadRequest, structs.ErrorRet{
			Error: "",
		})
	}

	err = integrator.AccountService.Register(accountadapter.RegisterReq{
		UserID:   req.UserID,
		Password: req.Password,
	})
	if err != nil {
		integrator.Logger.NewError(fmt.Sprintf("[Integrator][RegisterAccount] Error when executing account service register, err: %s", err.Error()))
		return c.JSON(http.StatusInternalServerError, structs.ErrorRet{
			Error: err.Error(),
		})
	}
	integrator.Logger.NewInfo(fmt.Sprintf("[Integrator][RegisterAccount] Account Registered, ID: %s", req.UserID))
	return c.JSON(http.StatusOK, nil)
}

type RegisterAccountReq struct {
	UserID   string `json:"userid"`
	Password string `json:"password"`
}
