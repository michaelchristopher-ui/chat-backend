package chatservice

import (
	"encoding/json"
	"log"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
)

// AddFriendReq is the json tagged request struct of the add friend feature
type AddFriendReq struct {
	FriendID string `json:"friend_id"`
}

// AddFriendResp is the json tagged response struct of the add friend feature
type AddFriendResp struct {
	Error string `json:"error"`
}

// addFriend is a function that adds the a<- is friends with ->b user relationship within the database
func (c ChatService) addFriend(userID string, data interface{}) {
	addFriendResp := AddFriendResp{}
	defer c.sendWebsocket(userID, addFriendResp)

	addFriendReq := AddFriendReq{}
	jsonString, err := json.Marshal(data)
	if json.Unmarshal(jsonString, &addFriendReq) != nil {
		log.Printf("fail unmarshal to publish message req")
	}
	if err != nil {
		addFriendResp.Error = err.Error()
		return
	}
	err = c.DB.AddFriend(databaseadapter.AddFriendReq{
		UserID:   userID,
		FriendID: addFriendReq.FriendID,
	})

	if err != nil {
		addFriendResp.Error = err.Error()
		log.Printf("error when adding friend: %s", err.Error())
		return
	}

}
