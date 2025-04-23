package wsprocessorservice

const (
	logUnrecognizedMessage = "websocket request type not recognized, request type %s"
	logFailUnmarshal       = "fail unmarshal %T to %T, err: %s"
	logFailMarshal         = "fail marshal %T, err: %s"
	logIncomingMessage     = "incoming message: %+v"
	logReadMessageErr      = "read message error: %s"

	errNilConnStr        = "connection is nil"
	errNilCallServiceStr = "call chat service function is nil"
	errCloseChan         = "Error when closing channel for user %s, error: %s"
)
