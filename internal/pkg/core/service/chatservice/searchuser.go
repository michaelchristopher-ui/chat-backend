package chatservice

import (
	"websocket_client/internal/pkg/core/adapter/chatadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
)

// SearchUser is a function that rudimentarily searches for a user through the use of a database query
func (c ChatService) SearchUser(data chatadapter.SearchUserReq) chatadapter.SearchUserResp {
	users, err := c.DB.SearchUser(databaseadapter.SearchUserRequest{
		UserID: data.SearchUserID,
	})
	if err != nil {
		c.Logger.NewError(err.Error())
	}
	userIDs := []string{}
	for _, user := range users {
		userIDs = append(userIDs, user.UserID)
	}
	return chatadapter.SearchUserResp{
		UserIDs: userIDs,
	}
}
