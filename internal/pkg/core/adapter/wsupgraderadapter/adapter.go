package wsupgraderadapter

import (
	"net/http"
	"websocket_client/internal/pkg/core/adapter/wsconnadapter"
)

//go:generate mockgen -source=adapter.go -package=wsupgraderadapter -destination=adapter_mock.go
type Adapter interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (wsconnadapter.Adapter, error)
}
