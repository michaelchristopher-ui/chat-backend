package chatservice

import "encoding/json"

type SearchFriendRequest struct {
	FriendID string `json:"friend_id"`
}

func (c ChatService) searchMessageHandler(userID string, data interface{}) {
	req := SearchFriendRequest{}
	reqByte, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(reqByte, &req)
	if err != nil {
		return
	}

}
