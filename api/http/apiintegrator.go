package http

import (
	"websocket_client/internal/conf"
	"websocket_client/internal/pkg/core/adapter/accountadapter"
	"websocket_client/internal/pkg/core/adapter/chatadapter"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"
	"websocket_client/internal/pkg/core/adapter/wsprocessoradapter"
	"websocket_client/internal/pkg/core/adapter/wsupgraderadapter"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// APIIntegrator is the struct for all API handler methods
type APIIntegrator struct {
	ChatService    chatadapter.Adapter
	AccountService accountadapter.Adapter
	Logger         loggeradapter.Adapter
	Upgrader       wsupgraderadapter.Adapter
	WsProcessor    wsprocessoradapter.Adapter
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
		Logger:         req.Logger,
		Upgrader:       req.Upgrader,
		WsProcessor:    req.WsProcessor,
	}
}

type NewAPIIntegratorReq struct {
	ChatService    chatadapter.Adapter
	AccountService accountadapter.Adapter
	Logger         loggeradapter.Adapter
	Upgrader       wsupgraderadapter.Adapter
	WsProcessor    wsprocessoradapter.Adapter
}

// API is a method that initializes the integrator and sets up all the APIs for the application
func API(req APIReq) {
	integrator := NewAPIIntegrator(NewAPIIntegratorReq{
		ChatService:    req.ChatService,
		AccountService: req.AccountService,
		Logger:         req.Logger,
		Upgrader:       req.Upgrader,
		WsProcessor:    req.WsProcessor,
	})

	chat := req.E.Group("")
	chat.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(conf.GetConfig().Server.RateLimit))))

	chat.POST("/receive", integrator.ReceiveMessage)
	chat.POST("/register", integrator.RegisterAccount)
	chat.GET("/health_check", integrator.HealthCheck)

	chatWs := req.E.Group("")
	chatWs.Use(integrator.Auth, integrator.AddContextID)
	chatWs.GET("/ws", integrator.HandleWebsocket)
}

type APIReq struct {
	E              *echo.Echo
	ChatService    chatadapter.Adapter
	AccountService accountadapter.Adapter
	Logger         loggeradapter.Adapter
	Upgrader       wsupgraderadapter.Adapter
	WsProcessor    wsprocessoradapter.Adapter
}
