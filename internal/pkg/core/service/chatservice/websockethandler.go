package chatservice

import (
	"encoding/json"
	"fmt"
	"log"
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
	// Attempt to upgrade the connection
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(req.ResponseWriter, req.Request, nil)
	if err != nil {
		log.Printf("[ChatService][WebsocketHandler] Error when upgrading for user %s, error: %s", req.UserID, err.Error())
		return err
	}

	// Defer responsible for handling errors and the closing of the connection.
	defer func() {
		if err != nil {
			log.Printf("[ChatService][WebsocketHandler] Error during websocket for user %s, error: %s", req.UserID, err.Error())
		}

		/*
			Closes the connection.
			The underlying implementation is a channel, so even if there is an error, the channel will be garbage collected after some time.
		*/
		err = conn.Close()
		if err != nil {
			log.Printf("[ChatService][WebsocketHandler] Error when closing channel for user %s, error: %s", req.UserID, err.Error())
		}

		// Delete user entry in the online users map
		delete(c.UserConnections, req.UserID)

		//Broadcast to friends that the user is offline
		c.broadcastEmptyMessageToFriends(req.UserID, typeUserOffline)
	}()

	// Add websocket connection entry to the map so that messages can be sent to the connected user.
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
		/*
			Generate FlowID for easy tracing
			The FlowID is a UUID. Do take note of the RPS to find out when we should reset or recycle the UUIDs.
		*/
		flowID := common.GenerateUUID()

		// Read sent message
		websocketRequest := WebsocketReq{}
		mt, msg, err := conn.ReadMessage()

		// Stop processing when an error occured or a request to close was sent
		if err != nil || mt == websocket.CloseMessage {
			break
		}

		// Unmarshals request and logs it
		err = json.Unmarshal(msg, &websocketRequest)
		if err != nil {
			c.Logger.NewError(fmt.Sprintf(logPrefix, "WebsocketHandler", fmt.Sprintf(logFailUnmarshal, websocketRequest), flowID))
			continue
		}
		c.Logger.NewInfo(fmt.Sprintf(logPrefix, "WebsocketHandler", fmt.Sprintf(logIncomingMessage, websocketRequest), flowID))

		/*
			Nearly all the features that a user can request for is done through websocket.
			The switch handles it.
		*/
		switch websocketRequest.ReqType {

		// Handles the request to send a message to another user
		case incomingMessageTypeMessage:
			c.sendMessageHandler(req.UserID, websocketRequest.Data)

		// Handles the add friend request
		case incomingMessageTypeAddFriend:
			c.addFriend(req.UserID, websocketRequest.Data)

		// Handles the Get Chat History between the current user and another user
		case incomingMessageTypeGetChatHistory:
			c.getChatHistory(req.UserID, websocketRequest.Data)

		// Handles the remove friend request
		case incomingMessageTypeRemoveFriend:
			c.removeFriend(req.UserID, websocketRequest.Data)

		// Handles the search friend request
		case incomingMessageTypeSearchFriend:
			c.removeFriend(req.UserID, websocketRequest.Data)

		// Default case logs invalid requests
		default:
			c.Logger.NewError(fmt.Sprintf(logPrefix, "WebsocketHandler", fmt.Sprintf(logUnrecognizedMessage, websocketRequest.ReqType), flowID))
		}
	}

	return nil
}
