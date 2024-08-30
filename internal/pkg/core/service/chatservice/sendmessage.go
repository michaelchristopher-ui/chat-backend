package chatservice

import (
	"encoding/json"
	"fmt"
	"time"
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
	publishMessageReq := PublishMessageReq{}
	jsonString, _ := json.Marshal(data)
	sendMessagesRes.Error = json.Unmarshal(jsonString, &publishMessageReq).Error()
	if sendMessagesRes.Error != "" {
		c.Logger.NewInfo(fmt.Sprintf(logFailUnmarshal, publishMessageReq))
		return
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
			c.Logger.NewError(fmt.Sprintf(logPrefix, "sendMessageHandler", fmt.Sprintf(logErrSaveMessage, err.Error()), "FlowIDTODO"))
			return err
		}

		err = c.publishMessage(publishMessageReq)
		if err != nil {
			c.Logger.NewError(fmt.Sprintf(logPrefix, "sendMessageHandler", fmt.Sprintf(logErrMessageCannotPublish, err.Error(), publishMessageReq), "FlowIDTODO"))
			return err
		}
		return nil
	}
	sendMessagesRes.Error = c.DB.DoCustomTransaction(transactionFunc).Error()
}
