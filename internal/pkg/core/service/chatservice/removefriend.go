package chatservice

import (
	"websocket_client/internal/pkg/core/adapter/chatadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
)

// removeFriend is a function that removes the a<- is friends with ->b user relationship within the database
func (c ChatService) RemoveFriend(removeFriendReq chatadapter.RemoveFriendReq) chatadapter.RemoveFriendResp {
	removeFriendResp := chatadapter.RemoveFriendResp{}

	err := c.DB.RemoveFriend(databaseadapter.RemoveFriendReq{
		UserID:   removeFriendReq.UserID,
		FriendID: removeFriendReq.FriendID,
	})

	if err != nil {
		removeFriendResp.Error = err.Error()
		c.Logger.NewError(LogErrRemoveMessage, err.Error())
		return removeFriendResp
	}
	c.Logger.NewInfo("success remove friend for userID %s", removeFriendReq.UserID)
	return removeFriendResp
}
