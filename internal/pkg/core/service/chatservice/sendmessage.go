package chatservice

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"websocket_client/internal/pkg/core/adapter/chatadapter"
	"websocket_client/internal/pkg/platform/mysql/models"

	"gorm.io/gorm"
)

// sendMessageRes is a json tagged response struct that defines the contents of the response of the handler
type SendMessageRes struct {
	Error string `json:"error"`
}

// sendMessageHandler is a function that saves and sends the specified message using a transaction
func (c ChatService) sendMessageHandler(userID string, data interface{}) {
	sendMessagesRes := SendMessageRes{}
	defer c.sendWebsocket(userID, sendMessagesRes)
	publishMessageReq := chatadapter.PublishMessageReq{}
	jsonString, _ := json.Marshal(data)
	if json.Unmarshal(jsonString, &publishMessageReq) != nil {
		log.Printf(logFailUnmarshal, sendMessagesRes)
	}

	currTimeSecond := time.Now().Second()

	publishMessageReq.FromUserID = userID

	publishMessageReq.Timestamp = fmt.Sprintf("%d", currTimeSecond)
	transactionFunc := func(tx *gorm.DB) error {
		err := tx.Create(models.Messages{
			Message:    publishMessageReq.Message,
			ToUserID:   publishMessageReq.ToUserID,
			FromUserID: publishMessageReq.FromUserID,
			Type:       typeMessage,
			Timestamp:  publishMessageReq.Timestamp,
		}).Error
		if err != nil {
			log.Printf(logSaveMessage, err.Error())
			return err
		}

		err = c.PublishMessage(publishMessageReq)
		if err != nil {
			log.Printf(logMessageCannotPublish, err.Error(), publishMessageReq)
			return err
		}
		return nil
	}
	sendMessagesRes.Error = c.DB.DoCustomTransaction(transactionFunc).Error()
}
