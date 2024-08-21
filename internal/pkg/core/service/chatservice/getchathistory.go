package chatservice

import (
	"encoding/json"
	"log"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
)

// GetChatHistoryRes is a struct that defines the contents of the websocket request message sent to the user after executing the getChatHistory function
type GetChatHistoryReq struct {
	FromUserID     string `json:"from_user_id"`
	ToUserID       string `json:"to_user_id"`
	Offset         int    `json:"offset"`
	Limit          int    `json:"limit"`
	TimestampAfter string `json:"timestamp_after"`
}

// GetChatHistoryRes is a struct that defines the contents of the websocket response message sent to the user after executing the getChatHistory function
type GetChatHistoryRes struct {
	Messages []Messages `json:"messages"`
	Error    string     `json:"error"`
}

// Messages is a struct that defines the contents of the returned messages
type Messages struct {
	FromUserID string `json:"from_user_id"`
	ToUserID   string `jsonL:"to_user_id"`
	Message    string `json:"message"`
	Type       int    `json:"type"`
	Timestamp  string `json:"timestamp"`
}

// getChatHistory is a function that retrieves chat history based on parameters supplied
func (c ChatService) getChatHistory(userID string, data interface{}) {
	getMessagesRes := GetChatHistoryRes{}
	defer c.sendWebsocket(userID, getMessagesRes)
	getMessagesReq := GetChatHistoryReq{}
	jsonString, _ := json.Marshal(data)
	if json.Unmarshal(jsonString, &getMessagesReq) != nil {
		log.Printf(logFailUnmarshal, getMessagesRes)
	}
	historyMessages, err := c.DB.GetChatHistory(databaseadapter.GetChatHistoryReq{
		FromUserID: getMessagesReq.FromUserID,
		ToUserID:   getMessagesReq.ToUserID,
		Offset:     getMessagesReq.Offset,
		Limit:      getMessagesReq.Limit,
	})
	if err != nil {
		getMessagesRes.Error = err.Error()
		return
	}

	for _, eachMessage := range historyMessages {
		getMessagesRes.Messages = append(getMessagesRes.Messages, Messages{
			FromUserID: eachMessage.FromUserID,
			ToUserID:   eachMessage.ToUserID,
			Message:    eachMessage.Message,
			Type:       eachMessage.Type,
			Timestamp:  eachMessage.Timestamp,
		})
	}

}
