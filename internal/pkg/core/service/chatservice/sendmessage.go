package chatservice

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"websocket_client/internal/common"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"

	"gorm.io/gorm"
)

// sendMessageRes is a json tagged response struct that defines the contents of the response of the handler
type SendMessageRes struct {
	Error string `json:"error"`
}

// sendMessageHandler is a function that saves and sends the specified message using a transaction
func (c ChatService) sendMessageHandler(userID string, data interface{}) {

	/*
		We want the sending user to be informed whether the sending of the message is a success
		Therefore we defer the act of informing the sending user first.
	*/
	sendMessagesRes := SendMessageRes{}
	defer c.sendWebsocket(userID, sendMessagesRes)

	/*
		Setup the payload by attempting to marshal the data parameter, then unmarshaling it again to the payload struct.
	*/
	publishMessageReq := PublishMessageReq{}
	jsonString, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error during marshaling of data: %s", err.Error())
		return
	}
	err = json.Unmarshal(jsonString, &publishMessageReq)
	if err != nil {
		sendMessagesRes.Error = err.Error()
		c.Logger.NewInfo(fmt.Sprintf(logFailUnmarshal, publishMessageReq))
		return
	}

	/*
		We get the current time, using it as a timestamp for the payload.
		In this case, we assume that the system clock is correct.
		In reality, we might have to have a centralized time server unless we can reasonably guarantee system clock correctness.
		For example, this could be through the use of a GPS clock.
	*/
	currTimeSecond := time.Now().Second()

	publishMessageReq.FromUserID = userID

	publishMessageReq.Timestamp = fmt.Sprintf("%d", currTimeSecond)
	transactionFunc := func(tx *gorm.DB) error {
		/*
			The message is saved to the database. The ID is an generated UUID.
			If the message cannot be saved we consider the sending of the message a failure.
		*/
		uuid := common.GenerateUUID()
		err := c.DB.SaveChatHistoryWithTx(tx, databaseadapter.SaveChatHistoryWithTxReq{
			ID:         uuid,
			Message:    publishMessageReq.Message,
			ToUserID:   publishMessageReq.ToUserID,
			FromUserID: publishMessageReq.FromUserID,
			Type:       typeMessage,
			Timestamp:  publishMessageReq.Timestamp,
		},
		)
		if err != nil {
			c.Logger.NewError(fmt.Sprintf(logPrefix, "[ChatService][sendMessageHandler]", fmt.Sprintf(logErrSaveMessage, err.Error()), "FlowIDTODO"))
			return err
		}

		/*
			At this time, we don't care if the user is online when we send the message.
			Should the user be offline, they will get the message through the fetching of the history.
		*/
		_, err = c.publishMessage(publishMessageReq)
		if err != nil {
			c.Logger.NewError(fmt.Sprintf(logPrefix, "[ChatService][sendMessageHandler]", fmt.Sprintf(logErrMessageCannotPublish, err.Error(), publishMessageReq), "FlowIDTODO"))
			return err
		}

		// Log on success
		c.Logger.NewInfo(fmt.Sprintf("[ChatService][sendMessageHandler] Success sending message with ID: %s", uuid))
		return nil
	}
	err = c.DB.DoCustomTransaction(transactionFunc)
	if err != nil {
		sendMessagesRes.Error = err.Error()
	}
}
