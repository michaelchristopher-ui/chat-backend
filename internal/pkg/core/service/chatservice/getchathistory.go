package chatservice

import (
	"websocket_client/internal/pkg/core/adapter/chatadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
)

// getChatHistory is a function that manually retrieves chat history based on parameters supplied
func (c ChatService) GetChatHistory(req chatadapter.GetChatHistoryReq) chatadapter.GetChatHistoryResp {
	getMessagesRes := chatadapter.GetChatHistoryResp{}

	historyMessages, err := c.DB.GetChatHistory(databaseadapter.GetChatHistoryReq{
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
		Offset:     req.Offset,
		Limit:      req.Limit,
	})
	if err != nil {
		getMessagesRes.Error = err.Error()
		return getMessagesRes
	}

	for _, eachMessage := range historyMessages {
		getMessagesRes.Messages = append(getMessagesRes.Messages, chatadapter.Message{
			FromUserID: eachMessage.FromUserID,
			ToUserID:   eachMessage.ToUserID,
			Message:    eachMessage.Message,
			Type:       "eachMessage.Type",
			Timestamp:  eachMessage.Timestamp,
		})
	}
	return getMessagesRes
}
