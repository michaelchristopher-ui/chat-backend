package chatservice

import (
	"websocket_client/internal/pkg/core/adapter/chatadapter"
)

/*
ReceiveMessage attempts to send the message to the user through websocket if they are connected to this server.
An error is returned when the user is online and the message fails to be sent through the websocket connection.
If the user is not online the server does not return an error, but isOnline will return false
*/
func (c ChatService) ReceiveMessage(req chatadapter.ReceiveMessageReq) (isOnline bool, err error) {
	conn, ok := c.UserConnections[req.ToUserID]
	if !ok {
		return false, nil
	}
	err = conn.WriteJSON(
		MessagePayload{
			FromUserID: req.FromUserID,
			Type:       req.Type,
			Message:    req.Message,
			Timestamp:  req.Timestamp,
			ToUserID:   req.ToUserID,
		},
	)

	return true, err
}
