package senderservice

import (
	"websocket_client/internal/common"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
	"websocket_client/internal/pkg/core/adapter/kvadapter"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"
	"websocket_client/internal/pkg/core/adapter/senderadapter"
)

type SenderService struct {
	DB     databaseadapter.RepoAdapter
	Redis  kvadapter.RepoAdapter
	Logger loggeradapter.Adapter
}

// NewSenderServiceReq defines the request parameter struct for NewSenderService
type NewSenderServiceReq struct {
	DB     databaseadapter.RepoAdapter
	Redis  kvadapter.RepoAdapter
	Logger loggeradapter.Adapter
}

// NewSenderService is a constructor function for SenderService, which conforms to senderadapter.Adapter
func NewSenderService(req NewSenderServiceReq) senderadapter.Adapter {
	if err := common.CheckNilFields(req); err != nil {
		panic(err.Error())
	}
	return SenderService(req)
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
	ToUserID   string `json:"to_user_id"`
}
