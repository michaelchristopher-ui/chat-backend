package http

//This group of constants define the possible incoming message types
const (
	incomingMessageTypeMessage        = "MESSAGE"
	incomingMessageTypeAddFriend      = "ADDFRIEND"
	incomingMessageTypeGetChatHistory = "GETCHATHISTORY"
	incomingMessageTypeRemoveFriend   = "REMOVEFRIEND"
	incomingMessageTypeSearchUser     = "SEARCHFRIEND"

	returnErrorPasswordLength = "Password should be shorter than 72 characters"
)
