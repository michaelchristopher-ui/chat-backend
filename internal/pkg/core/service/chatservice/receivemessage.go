package chatservice

/*
ReceiveMessage is a function that:

1. TODO: Lock the user connection if it is still used?

2. Attempts to send a message to the user's websocket if the connection is available

3. Keep retrying when it fails to send the message as long as the user is connected
*/

func (c ChatService) ReceiveMessage(userID string, data interface{}) (isOnline bool, err error) {
	conn := c.WsStore.GetConn(userID)
	if conn != nil {
		defer c.WsStore.Release(userID)
		err = conn.WriteJSON(data)
		return true, err
	}
	return false, nil
}
