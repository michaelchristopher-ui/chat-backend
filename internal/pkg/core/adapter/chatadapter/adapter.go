package chatadapter

import "net/http"

// Adapter defines an interface for the chat feature
type Adapter interface {
	WebsocketHandler(req WebsocketHandlerReq) error
	PublishMessage(req PublishMessageReq) error
	ReceiveMessage(req ReceiveMessageReq) error
}

type AddFriendReq struct {
	UserID   string
	FriendID string
}

type WebsocketHandlerReq struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	UserID         string
}

type PublishMessageReq struct {
	Message    string `json:"message"`
	Type       int    `json:"type"`
	ToUserID   string `json:"to_user_id"`
	FromUserID string `json:"-"`
	Timestamp  string `json:"-"`
}

type ReceiveMessageReq struct {
	Message    string
	FromUserID string
	Type       int
	ToUserID   string
	Timestamp  string
}

type ReceiveMessageRes struct{}

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
