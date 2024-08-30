package chatservice

import (
	"encoding/json"
	"fmt"
	"websocket_client/internal/common"
	"websocket_client/internal/pkg/core/adapter/chatadapter"

	"github.com/gorilla/websocket"
)

// WebsocketReq is a struct that defines the contents of incoming websocket messages
type WebsocketReq struct {
	ReqType string      `json:"req_type"`
	Data    interface{} `json:"data"`
}

/*
WebsocketHandler is a function that handles the upgrading of the connection and the processing of messages
while also managing an entry in Redis so messages can be sent to the user.
*/
func (c ChatService) WebsocketHandler(req chatadapter.WebsocketHandlerReq) error {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(req.ResponseWriter, req.Request, nil)
	if err != nil {
		return err
	}
	defer func() {
		conn.Close()
		delete(c.UserConnections, req.UserID)
		//Broadcast to friends that the user is offline
		c.broadcastEmptyMessageToFriends(req.UserID, typeUserOffline)
	}()
	c.UserConnections[req.UserID] = conn

	//Broadcast to friends that the user is online
	c.broadcastEmptyMessageToFriends(req.UserID, typeUserOnline)

	//Setup user and connected server key-value entry to redis so that services can find out which server the user is connected to
	isOpen := true
	c.Redis.SetValueUntilChannelClose(req.UserID, *common.IPPort, 30, &isOpen)
	defer func() {
		isOpen = false
	}()

	//Process incoming websocket messages
	for {
		flowID := common.GenerateUUID()
		websocketRequest := WebsocketReq{}
		mt, msg, err := conn.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}
		err = json.Unmarshal(msg, &websocketRequest)
		if err != nil {
			c.Logger.NewError(fmt.Sprintf(logPrefix, "WebsocketHandler", fmt.Sprintf(logFailUnmarshal, websocketRequest), flowID))
			continue
		}
		c.Logger.NewInfo(fmt.Sprintf(logPrefix, "WebsocketHandler", fmt.Sprintf(logIncomingMessage, websocketRequest), flowID))
		switch websocketRequest.ReqType {
		case incomingMessageTypeMessage:
			c.sendMessageHandler(req.UserID, websocketRequest.Data)
		case incomingMessageTypeAddFriend:
			c.addFriend(req.UserID, websocketRequest.Data)
		case incomingMessageTypeGetChatHistory:
			c.getChatHistory(req.UserID, websocketRequest.Data)
		case incomingMessageTypeRemoveFriend:
			c.removeFriend(req.UserID, websocketRequest.Data)
		default:
			c.Logger.NewError(fmt.Sprintf(logPrefix, "WebsocketHandler", fmt.Sprintf(logUnrecognizedMessage, websocketRequest.ReqType), flowID))
		}
	}

	return nil
}
