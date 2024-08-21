package chatservice

import (
	"websocket_client/internal/pkg/core/adapter/chatadapter"
)

/*
ReceiveMessage attempts to send the message to the user through websocket if they are connected to this server.
An error is returned when the user is online and the message fails to be sent through the websocket connection.
If the user is not online the server does not return an error.
*/
func (c ChatService) ReceiveMessage(req chatadapter.ReceiveMessageReq) error {
	if conn, ok := c.UserConnections[req.ToUserID]; ok {
		return conn.WriteJSON(
			MessagePayload{
				FromUserID: req.FromUserID,
				Type:       req.Type,
				Message:    req.Message,
				Timestamp:  req.Timestamp,
			},
		)
	}

	return nil
}
