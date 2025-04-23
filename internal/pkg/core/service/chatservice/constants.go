package chatservice

//This group of constants define the possible log message formats
const (
	LogErrMessageCannotPublish = "message cannot be published, err: %s, message data: %+v"
	LogErrAddMessage           = "error when adding friend: %s"
	LogErrRemoveMessage        = "error when removing friend: %s"
	LogErrSaveMessage          = "error when saving message: %s"
	LogErrAddFriend            = "error when adding friend for user id %s: %s"
)

//This group of constants define all the possible message types
const (
	typeUserOnline = iota
	typeMessage
	typeUserOffline
)
