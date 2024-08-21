package http

import (
	"encoding/json"
	"net/http"
	"websocket_client/api/http/structs"
	"websocket_client/internal/pkg/core/adapter/chatadapter"

	"github.com/labstack/echo"
)

// ReceiveMessage is the handler method for the /receive api endpoint.
func (integrator *APIIntegrator) ReceiveMessage(c echo.Context) error {
	req := ReceiveMessageReq{}
	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, structs.ErrorRet{
			Error: err.Error(),
		})
	}
	err = integrator.ChatService.ReceiveMessage(chatadapter.ReceiveMessageReq{
		Message:    req.Message,
		FromUserID: req.FromUserID,
		Type:       req.Type,
		ToUserID:   req.ToUserID,
		Timestamp:  req.Timestamp,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, structs.ErrorRet{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, nil)
}

type ReceiveMessageReq struct {
	Message    string `json:"message"`
	FromUserID string `json:"from_user_id"`
	Type       int    `json:"type"`
	ToUserID   string `json:"to_user_id"`
	Timestamp  string `json:"timestamp"`
}
