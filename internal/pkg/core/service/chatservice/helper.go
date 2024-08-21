package chatservice

import (
	"log"
	"websocket_client/internal/pkg/core/adapter/chatadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
)

/*
sendWebsocket is a function that:
1. Attempts to send a message to the user's websocket if the connection is available
2. Keep retrying when it fails to send the message as long as the user is connected
*/
func (c ChatService) sendWebsocket(userID string, data interface{}) {
	if conn, ok := c.UserConnections[userID]; ok {
		var err error
		for ok {
			err = conn.WriteJSON(data)
			if err == nil {
				break
			}
			conn, ok = c.UserConnections[userID]
		}
	}
}

// broadcastEmptyMessageToFriends is a function intended to send an empty message to the user's friends indicating the user's online/offline status
func (c ChatService) broadcastEmptyMessageToFriends(userID string, messageType int) {
	userFriends, err := c.DB.GetUserFriends(databaseadapter.GetUserFriendsReq{
		UserID: userID,
	})

	if err != nil {
		for _, eachUserFriend := range userFriends {
			publishMessageReq := chatadapter.PublishMessageReq{
				Message:    "",
				ToUserID:   eachUserFriend.UserFriendID,
				FromUserID: userID,
				Type:       messageType,
			}
			err := c.PublishMessage(publishMessageReq)
			if err != nil {
				log.Printf(logMessageCannotPublish, err.Error(), publishMessageReq)
			}
		}
	}
}
