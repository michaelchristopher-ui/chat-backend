package chatservice

import (
	"websocket_client/internal/common"
	"websocket_client/internal/pkg/core/adapter/chatadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
	"websocket_client/internal/pkg/core/adapter/kvadapter"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"

	"github.com/gorilla/websocket"
)

// ChatService is a service that handles all chat-related functions
type ChatService struct {
	DB              databaseadapter.RepoAdapter
	UserConnections map[string]*websocket.Conn
	Redis           kvadapter.RepoAdapter
	Logger          loggeradapter.Adapter
}

// NewChatServiceReq defines the request parameter struct for NewChatService
type NewChatServiceReq struct {
	DB     databaseadapter.RepoAdapter
	Redis  kvadapter.RepoAdapter
	Logger loggeradapter.Adapter
}

// NewChatService is a constructor function for ChatService, which conforms to chatadapter.Adapter
func NewChatService(req NewChatServiceReq) chatadapter.Adapter {
	if err := common.CheckNilFields(req); err != nil {
		panic(err.Error())
	}
	return ChatService{
		DB:              req.DB,
		UserConnections: map[string]*websocket.Conn{},
		Redis:           req.Redis,
		Logger:          req.Logger,
	}
}

/*
MessagePayload is a type that specifies the json contents of the message that is to be sent by the PublishMessage function
and also received by the recipient in the ReceiveMessage function
*/
type MessagePayload struct {
	Message    string `json:"message"`
	Type       int    `json:"type"`
	FromUserID string `json:"from_user_id"`
	Timestamp  string `json:"timestamp"`
}
