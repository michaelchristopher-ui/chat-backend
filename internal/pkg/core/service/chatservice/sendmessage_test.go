package chatservice_test

import (
	"errors"
	"testing"
	"websocket_client/internal/pkg/core/adapter/chatadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"
	"websocket_client/internal/pkg/core/adapter/senderadapter"
	"websocket_client/internal/pkg/core/service/chatservice"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSendMessageTransactionFuncFactory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := databaseadapter.NewMockRepoAdapter(ctrl)
	mockSender := senderadapter.NewMockAdapter(ctrl)
	mockLogger := loggeradapter.NewMockAdapter(ctrl)

	service := chatservice.NewChatService(chatservice.NewChatServiceReq{
		DB:     mockDB,
		Sender: mockSender,
		Logger: mockLogger,
	})

	mockDB.EXPECT().SaveChatHistoryWithTx(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockSender.EXPECT().PublishMessage(gomock.Any()).Return(true, nil).Times(1)
	mockLogger.EXPECT().NewInfo(gomock.Any(), gomock.Any()).Times(1)

	sendMessageTransactionFunc := service.SendMessageTransactionFuncFactory(chatadapter.SendMessageReq{})
	err := sendMessageTransactionFunc(nil)

	assert.Nil(t, err)
}

func TestSendChatService_SendMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := databaseadapter.NewMockRepoAdapter(ctrl)

	service := chatservice.NewChatService(chatservice.NewChatServiceReq{
		DB: mockDB,
	})

	type SendMessagesTestData struct {
		name string
		req  chatadapter.SendMessageReq
		mock func()
	}

	tests := []SendMessagesTestData{
		{
			name: "Test success sending message",
			mock: func() {
				mockDB.EXPECT().DoCustomTransaction(gomock.Any()).Return(nil).Times(1)
			},
			req: chatadapter.SendMessageReq{
				FromUserID: "user1",
				ToUserID:   "user2",
				Message:    "Hello",
				Type:       1,
				Timestamp:  "2025-04-02T06:53:00Z",
			},
		}, {
			name: "Test fail sending message",
			mock: func() {
				mockDB.EXPECT().DoCustomTransaction(gomock.Any()).Return(errors.New("foo")).Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			service.SendMessage(
				test.req,
			)
		})
	}

}
