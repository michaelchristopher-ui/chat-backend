package chatservice

import (
	"websocket_client/internal/pkg/core/adapter/chatadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
)

// addFriend is a function that adds the a<- is friends with ->b user relationship within the database
func (c ChatService) AddFriend(addFriendReq chatadapter.AddFriendReq) chatadapter.AddFriendResp {
	addFriendResp := chatadapter.AddFriendResp{}

	err := c.DB.AddFriend(databaseadapter.AddFriendReq{
		UserID:   addFriendReq.UserID,
		FriendID: addFriendReq.FriendID,
	})

	if err != nil {
		addFriendResp.Error = err.Error()
		c.Logger.NewError(LogErrAddFriend, addFriendReq.UserID, err.Error())
		return addFriendResp
	}
	return addFriendResp
}
