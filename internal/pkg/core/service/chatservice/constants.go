package chatservice

//This group of constants define the possible log message formats
const (
	logFailUnmarshal        = "fail unmarshal to  %T"
	logIPEmpty              = "ip is empty for %s"
	logMessageCannotPublish = "message cannot be published, err: %s, message data: %+v"
	logAddMessage           = "error when adding friend: %s"
	logRemoveMessage        = "error when removing friend: %s"
	logSaveMessage          = "error when saving message: %s"
	logUnrecognizedMessage  = "websocket request not recognized"
)

//This group of constants define thhe possible incoming message types
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
