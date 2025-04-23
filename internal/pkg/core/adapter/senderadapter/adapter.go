package senderadapter

//go:generate mockgen -source=adapter.go -package=senderadapter -destination=adapter_mock.go
type Adapter interface {
	PublishMessage(req SendMessageReq) (isOnline bool, err error)
}

type SendMessageReq struct {
	ID         string
	Message    string
	Type       int
	ToUserID   string
	FromUserID string
	Timestamp  string
}
