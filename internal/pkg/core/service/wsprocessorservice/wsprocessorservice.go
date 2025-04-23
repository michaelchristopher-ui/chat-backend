package wsprocessorservice

import (
	"encoding/json"
	"errors"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"
	"websocket_client/internal/pkg/core/adapter/wsconnadapter"
	"websocket_client/internal/pkg/core/adapter/wsprocessoradapter"
	"websocket_client/internal/pkg/core/adapter/wsstoreadapter"

	"github.com/gorilla/websocket"
)

type WsProcessorService struct {
	Logger  loggeradapter.Adapter
	WsStore wsstoreadapter.Adapter
}

type WsProcessorServiceReq struct {
	Logger  loggeradapter.Adapter
	WsStore wsstoreadapter.Adapter
}

func NewWsProcessorService(req WsProcessorServiceReq) wsprocessoradapter.Adapter {
	return &WsProcessorService{
		Logger:  req.Logger,
		WsStore: req.WsStore,
	}
}

/*
ProcessWebsocketAfterUpgrade accepts a websocket connecction, userID (derived from the context in the websocket handler) and a callChatServiceFunc.
It reads incoming messages sent to the websocket connection, unmarshals the message, separates the byte array data and the request type,
and calls the callChatServiceFunc, supplying the extracted byte array, the userID and the request type as parameters.

This function is designed to be called after the upgrader successfully upgrades the connection
*/
func (p *WsProcessorService) ProcessWebsocketAfterUpgrade(conn wsconnadapter.Adapter, userID string, callChatServiceFunc func(dataByte []byte, userID, reqType string) error) (err error) {
	//Check if required nullable paramms are not nil (concrete type and value) before processing
	if conn == nil {
		return errors.New(errNilConnStr)
	}
	if callChatServiceFunc == nil {
		return errors.New(errNilCallServiceStr)
	}

	defer func() {
		//Deletes the entry in the map
		p.WsStore.DeleteConn(userID)

		/*
			Closes the connection.
			The underlying implementation is a channel, so even if there is an error, the channel will be garbage collected after some time.
		*/
		err = conn.Close()
		if err != nil {
			p.Logger.NewError(errCloseChan, userID, err.Error())
		}
	}()
	err = p.WsStore.AddConn(conn, userID)
	if err != nil {
		return err
	}
	//Process incoming websocket messages
	for {
		// Read sent message, continue if there are no errors
		websocketRequest := WebsocketReq{}
		var mt int
		var msg []byte
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			p.Logger.NewError(logReadMessageErr, err.Error())
		}

		// Stop processing when an error occured or a request to close was sent
		if err != nil || mt == websocket.CloseMessage {
			break
		}

		// Unmarshals request and logs it
		err = json.Unmarshal(msg, &websocketRequest)
		if err != nil {
			p.Logger.NewError(logFailUnmarshal, msg, websocketRequest, err.Error())
			continue
		}
		p.Logger.NewInfo(logIncomingMessage, websocketRequest)

		//Marshal the data parameter for chat service func
		dataByte, err := json.Marshal(websocketRequest.Data)
		if err != nil {
			p.Logger.NewError(logFailMarshal, websocketRequest.Data, err.Error())
			continue
		}

		//Calls the Chat Service class methods depending on the request type
		err = callChatServiceFunc(dataByte, userID, websocketRequest.ReqType)
		if err != nil {
			p.Logger.NewError(logFailMarshal, websocketRequest.Data, err.Error())
			continue
		}

	}

	return nil
}

// WebsocketReq is a struct that defines the contents of incoming websocket messages
type WebsocketReq struct {
	ReqType string      `json:"req_type"`
	Data    interface{} `json:"data"`
}
