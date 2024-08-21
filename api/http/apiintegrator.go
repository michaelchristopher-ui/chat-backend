package http

import (
	"websocket_client/internal/pkg/core/adapter/accountadapter"
	"websocket_client/internal/pkg/core/adapter/chatadapter"

	"github.com/labstack/echo"
)

// APIIntegrator is the struct for all API handler methods
type APIIntegrator struct {
	ChatService    chatadapter.Adapter
	AccountService accountadapter.Adapter
}

/*
NewAPIIntegrator creates a new APIIntegrator instance,
sets up all the required service components
and returns a pointer to it
*/
func NewAPIIntegrator(req NewAPIIntegratorReq) *APIIntegrator {
	return &APIIntegrator{
		ChatService:    req.ChatService,
		AccountService: req.AccountService,
	}
}

type NewAPIIntegratorReq struct {
	ChatService    chatadapter.Adapter
	AccountService accountadapter.Adapter
}

// API is a method that initializes the integrator and sets up all the APIs for the application
func API(req APIReq) {
	integrator := NewAPIIntegrator(NewAPIIntegratorReq{
		ChatService:    req.ChatService,
		AccountService: req.AccountService,
	})

	chat := req.E.Group("")

	chat.GET("/ws", integrator.Websocket)
	chat.POST("/receive", integrator.ReceiveMessage)
	chat.POST("/register", integrator.RegisterAccount)
}

type APIReq struct {
	E              *echo.Echo
	ChatService    chatadapter.Adapter
	AccountService accountadapter.Adapter
}
