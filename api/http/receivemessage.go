package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"websocket_client/api/http/structs"
	"websocket_client/internal/pkg/core/adapter/chatadapter"

	"github.com/labstack/echo/v4"
)

// ReceiveMessage is the handler method for the /receive api endpoint.
func (integrator *APIIntegrator) ReceiveMessage(c echo.Context) error {
	req := ReceiveMessageReq{}
	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		integrator.Logger.NewError(fmt.Sprintf("[Integrator][ReceiveMesssage] Error when decoding request, err: %s", err.Error()))
		return c.JSON(http.StatusInternalServerError, structs.ErrorRet{
			Error: err.Error(),
		})
	}
	isOnline, err := integrator.ChatService.ReceiveMessage(req.ToUserID, chatadapter.Message{
		Message:    req.Message,
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
		Timestamp:  req.Timestamp,
	})
	if err != nil {
		integrator.Logger.NewError(fmt.Sprintf("[Integrator][ReceiveMesssage] Error when sending message to user with ID %s, online: %v, err: %s", req.ToUserID, isOnline, err.Error()))
		return c.JSON(http.StatusInternalServerError, structs.ErrorRet{
			Error: err.Error(),
		})
	}

	integrator.Logger.NewInfo(fmt.Sprintf("[Integrator][ReceiveMesssage] Received message from user with ID %s to ID %s, user online = %v", req.FromUserID, req.ToUserID, isOnline))
	return c.JSON(http.StatusOK, ReceiveMessageRes{
		IsOnline: isOnline,
	})
}

type ReceiveMessageReq struct {
	Message    string `json:"message"`
	FromUserID string `json:"from_user_id"`
	Type       int    `json:"type"`
	ToUserID   string `json:"to_user_id"`
	Timestamp  string `json:"timestamp"`
}

type ReceiveMessageRes struct {
	IsOnline bool `json:"is_online"`
}
