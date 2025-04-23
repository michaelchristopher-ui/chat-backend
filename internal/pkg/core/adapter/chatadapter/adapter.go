package chatadapter

import (
	"net/http"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// Adapter defines an interface for the chat feature
//
//go:generate mockgen -source=adapter.go -package=chatadapter -destination=adapter_mock.go
type Adapter interface {
	ReceiveMessage(userID string, message any) (isOnline bool, err error)
	SendMessage(data SendMessageReq) SendMessageResp
	SearchUser(data SearchUserReq) SearchUserResp
	RemoveFriend(data RemoveFriendReq) RemoveFriendResp
	GetChatHistory(data GetChatHistoryReq) GetChatHistoryResp
	AddFriend(data AddFriendReq) AddFriendResp
	SendMessageTransactionFuncFactory(publishMessageReq SendMessageReq) func(tx *gorm.DB) error
}

// sendMessageRes is a json tagged response struct that defines the contents of the response of the handler
type SendMessageResp struct {
	Error string `json:"error"`
}

type SearchUserResp struct {
	UserIDs []string `json:"user_ids"`
}

type GetChatHistoryReq struct {
	FromUserID     string
	ToUserID       string
	Offset         int
	Limit          int
	TimestampAfter string
}

// GetChatHistoryResp is a struct that defines the contents of the websocket response message sent to the user after executing the getChatHistory function
type GetChatHistoryResp struct {
	Messages []Message `json:"messages"`
	Error    string    `json:"error"`
}

// RemoveFriendResp is the json tagged response struct of the remove friend feature
type RemoveFriendResp struct {
	Error string `json:"error"`
}

type AddFriendReq struct {
	UserID   string
	FriendID string
}

type WebsocketHandlerReq struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	UserID         string
	Conn           *websocket.Conn
}

type ReceiveMessageReq struct {
	Message    string
	FromUserID string
	Type       int
	ToUserID   string
	Timestamp  string
}

type GetMessagesReq struct {
	FromUserID     string
	ToUserID       string
	Offset         int
	Limit          int
	TimestampAfter string
}

type GetMessagesRes struct {
	Messages []Message
}

type Message struct {
	FromUserID string
	ToUserID   string
	Message    string
	Type       string
	Timestamp  string
}

type AddFriendResp struct {
	Error string
}

type SendMessageReq struct {
	ID         string
	Message    string
	Type       int
	ToUserID   string
	FromUserID string
	Timestamp  string
}

type RemoveFriendReq struct {
	UserID   string
	FriendID string
}

type SearchUserReq struct {
	SearchUserID string
}
