package wsupgrader

import (
	"net/http"
	"websocket_client/internal/pkg/core/adapter/wsconnadapter"
	"websocket_client/internal/pkg/core/adapter/wsupgraderadapter"

	"github.com/gorilla/websocket"
)

type NewWsUpgraderServiceReq struct {
	Upgrader websocket.Upgrader
}

type WsUpgraderService struct {
	upgrader websocket.Upgrader
}

func NewWsUpgraderService(req NewWsUpgraderServiceReq) wsupgraderadapter.Adapter {
	return &WsUpgraderService{
		upgrader: req.Upgrader,
	}
}

func (ws *WsUpgraderService) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (wsconnadapter.Adapter, error) {
	return ws.upgrader.Upgrade(w, r, responseHeader)
}
