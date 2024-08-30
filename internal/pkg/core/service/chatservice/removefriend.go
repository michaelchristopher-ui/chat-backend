package chatservice

import (
	"encoding/json"
	"log"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
)

// RemoveFriendReq is the json tagged request struct of the remove friend feature
type RemoveFriendReq struct {
	FriendID string `json:"friend_id"`
}

// RemoveFriendResp is the json tagged response struct of the remove friend feature
type RemoveFriendResp struct {
	Error string `json:"error"`
}

// removeFriend is a function that removes the a<- is friends with ->b user relationship within the database
func (c ChatService) removeFriend(userID string, data interface{}) {
	removeFriendResp := RemoveFriendResp{}
	defer c.sendWebsocket(userID, removeFriendResp)

	removeFriendReq := RemoveFriendReq{}
	jsonString, err := json.Marshal(data)
	if json.Unmarshal(jsonString, &removeFriendReq) != nil {
		log.Printf(logFailUnmarshal, removeFriendReq)
	}
	if err != nil {
		removeFriendResp.Error = err.Error()
		return
	}
	err = c.DB.RemoveFriend(databaseadapter.RemoveFriendReq{
		UserID:   userID,
		FriendID: removeFriendReq.FriendID,
	})

	if err != nil {
		removeFriendResp.Error = err.Error()
		log.Printf(logErrRemoveMessage, err.Error())
		return
	}
}
