package wsprocessoradapter

import "websocket_client/internal/pkg/core/adapter/wsconnadapter"

//go:generate mockgen -source=adapter.go -package=wsprocessoradapter -destination=adapter_mock.go
type Adapter interface {
	ProcessWebsocketAfterUpgrade(conn wsconnadapter.Adapter, userID string, callChatServiceFunc func(dataByte []byte, userID, reqType string) error) (err error)
}
