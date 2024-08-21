package http

import (
	"log"
	"net/http"
	"websocket_client/api/http/structs"
	"websocket_client/internal/common"
	"websocket_client/internal/pkg/core/adapter/accountadapter"
	"websocket_client/internal/pkg/core/adapter/chatadapter"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

// Websocket is the handler method for the /ws api endpoint
func (integrator *APIIntegrator) Websocket(c echo.Context) error {
	userId, password, err := common.SplitUserIDAndPasswordFromAuth(c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, structs.ErrorRet{
			Error: err.Error(),
		})
	}

	err = integrator.AccountService.VerifyAuth(accountadapter.VerifyAuthReq{
		UserID:   userId,
		Password: password,
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, structs.ErrorRet{
			Error: err.Error(),
		})
	}
	log.Printf("userid: %s", userId)
	req := chatadapter.WebsocketHandlerReq{
		ResponseWriter: c.Response(),
		Request:        c.Request(),
		UserID:         userId,
	}
	err = integrator.ChatService.WebsocketHandler(req)
	if err != nil {
		return c.JSON(websocket.CloseAbnormalClosure, structs.ErrorRet{
			Error: err.Error(),
		})
	}
	return c.JSON(websocket.CloseNormalClosure, nil)
}
