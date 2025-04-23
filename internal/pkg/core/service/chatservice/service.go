package chatservice

import (
	"websocket_client/internal/pkg/core/adapter/chatadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
	"websocket_client/internal/pkg/core/adapter/kvadapter"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"
	"websocket_client/internal/pkg/core/adapter/senderadapter"
	"websocket_client/internal/pkg/core/adapter/wsstoreadapter"
)

// ChatService is a service that handles all chat-related functions
type ChatService struct {
	DB      databaseadapter.RepoAdapter
	Redis   kvadapter.RepoAdapter
	Logger  loggeradapter.Adapter
	Sender  senderadapter.Adapter
	WsStore wsstoreadapter.Adapter
}

// NewChatServiceReq defines the request parameter struct for NewChatService
type NewChatServiceReq struct {
	DB      databaseadapter.RepoAdapter
	Redis   kvadapter.RepoAdapter
	Logger  loggeradapter.Adapter
	Sender  senderadapter.Adapter
	WsStore wsstoreadapter.Adapter
}

// NewChatService is a constructor function for ChatService, which conforms to chatadapter.Adapter
func NewChatService(req NewChatServiceReq) chatadapter.Adapter {
	return ChatService(req)
}
