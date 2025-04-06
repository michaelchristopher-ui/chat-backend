package chatadapter

import "net/http"

//go:generate mockgen -source=adapter.go -package=chatadapter -destination=adapter_mock.go
// Adapter defines an interface for the chat feature
type Adapter interface {
	WebsocketHandler(req WebsocketHandlerReq) error
	ReceiveMessage(req ReceiveMessageReq) (bool, error)
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
