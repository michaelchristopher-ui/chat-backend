package chatservice

//This group of constants define the possible log message formats
const (
	logPrefix                  = "[ChatService][%s] %s FlowID: %s"
	logFailUnmarshal           = "fail unmarshal to  %T"
	logIPEmpty                 = "ip is empty for %s"
	logErrMessageCannotPublish = "message cannot be published, err: %s, message data: %+v"
	logErrAddMessage           = "error when adding friend: %s"
	logErrRemoveMessage        = "error when removing friend: %s"
	logErrSaveMessage          = "error when saving message: %s"
	logUnrecognizedMessage     = "websocket request type not recognized, request type %s"
	logIncomingMessage         = "incoming message: %+v"
	logErrDecode               = "error when decoding, %s"
)

//This group of constants define the possible incoming message types
const (
	incomingMessageTypeMessage        = "MESSAGE"
	incomingMessageTypeAddFriend      = "ADDFRIEND"
	incomingMessageTypeGetChatHistory = "GETCHATHISTORY"
	incomingMessageTypeRemoveFriend   = "REMOVEFRIEND"
)

//This group of constants define all the possible message types
const (
	typeUserOnline = iota
	typeMessage
	typeUserOffline
)
