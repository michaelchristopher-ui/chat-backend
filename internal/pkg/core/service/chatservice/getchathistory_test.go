package chatservice_test

import (
	"errors"
	"testing"
	"websocket_client/internal/pkg/core/adapter/chatadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
	"websocket_client/internal/pkg/core/adapter/kvadapter"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"
	"websocket_client/internal/pkg/core/service/chatservice"
	"websocket_client/internal/pkg/platform/mysql/models"

	"github.com/golang/mock/gomock"
)

func TestGetChatHistory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := loggeradapter.NewMockAdapter(mockCtrl)
	mockRedisClient := kvadapter.NewMockRepoAdapter(mockCtrl)
	mockDB := databaseadapter.NewMockRepoAdapter(mockCtrl)

	c := chatservice.NewChatService(chatservice.NewChatServiceReq{
		DB:     mockDB,
		Logger: mockLogger,
		Redis:  mockRedisClient,
	})

	type GetChatHistoryTestData struct {
		name string
		req  chatadapter.GetChatHistoryReq
		mock func()
	}

	tests := []GetChatHistoryTestData{
		{
			name: "Successful Retrieval",
			req: chatadapter.GetChatHistoryReq{
				FromUserID: "friend456",
				ToUserID:   "user123",
				Offset:     0,
				Limit:      10,
			},
			mock: func() {
				// Mocking successful DB call returning messages
				mockDB.EXPECT().GetChatHistory(gomock.Any()).Return([]models.Messages{
					{FromUserID: "friend456", ToUserID: "user123", Message: "Hello!"},
					{FromUserID: "user123", ToUserID: "friend456", Message: "Hi there!"},
				}, nil).Times(1)
			},
		},
		{
			name: "Database Error",
			req: chatadapter.GetChatHistoryReq{
				FromUserID: "friend456",
				ToUserID:   "user123",
				Offset:     0,
				Limit:      10,
			},
			mock: func() {
				// Mocking DB call returning an error
				mockDB.EXPECT().GetChatHistory(gomock.Any()).Return(nil, errors.New("database connection failed")).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock() // Call mocked functions

			c.GetChatHistory(tt.req) // Call GetchatHistory method
		})
	}
}
