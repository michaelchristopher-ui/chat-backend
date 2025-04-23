package wsstoreadapter

import "websocket_client/internal/pkg/core/adapter/wsconnadapter"

//go:generate mockgen -source=adapter.go -package=wsstoreadapter -destination=adapter_mock.go
type Adapter interface {
	AddConn(conn wsconnadapter.Adapter, userID string) error
	GetConn(userID string) wsconnadapter.Adapter
	DeleteConn(userID string)
	Release(userID string)
}
