package chatservice

import (
	"fmt"
	"time"
	"websocket_client/internal/common"
	"websocket_client/internal/pkg/core/adapter/chatadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
	"websocket_client/internal/pkg/core/adapter/senderadapter"

	"gorm.io/gorm"
)

func (c ChatService) SendMessageTransactionFuncFactory(publishMessageReq chatadapter.SendMessageReq) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		/*
			The message is saved to the database. The ID is a generated UUID.
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
			c.Logger.NewError(LogErrSaveMessage, err.Error())
			return err
		}

		/*
			At this time, we don't care if the user is online when we send the message.
			Should the user be offline, they will get the message through the fetching of the history.
		*/
		_, err = c.Sender.PublishMessage(senderadapter.SendMessageReq{
			ID:         uuid,
			Message:    publishMessageReq.Message,
			ToUserID:   publishMessageReq.ToUserID,
			FromUserID: publishMessageReq.FromUserID,
			Type:       typeMessage,
			Timestamp:  publishMessageReq.Timestamp,
		})
		if err != nil {
			c.Logger.NewError(LogErrMessageCannotPublish, err.Error(), publishMessageReq)
			return err
		}

		// Log on success
		c.Logger.NewInfo(fmt.Sprintf("[ChatService][sendMessageHandler] Success sending message with ID: %s", uuid))
		return nil
	}
}

// sendMessage is a function that saves and sends the specified message using a transaction
func (c ChatService) SendMessage(publishMessageReq chatadapter.SendMessageReq) chatadapter.SendMessageResp {
	/*
		We want the sending user to be informed whether the sending of the message is a success
		Therefore we defer the act of informing the sending user first.
	*/
	sendMessagesRes := chatadapter.SendMessageResp{}
	/*
		We get the current time, using it as a timestamp for the payload.
		In this case, we assume that the system clock is correct.
		In reality, we might have to have a centralized time server unless we can reasonably guarantee system clock correctness.
		For example, this could be through the use of a GPS clock.
	*/
	currTimeSecond := time.Now().Second()

	publishMessageReq.Timestamp = fmt.Sprintf("%d", currTimeSecond)
	transactionFunc := c.SendMessageTransactionFuncFactory(publishMessageReq)
	err := c.DB.DoCustomTransaction(transactionFunc)
	if err != nil {
		sendMessagesRes.Error = err.Error()
	}
	return sendMessagesRes
}
