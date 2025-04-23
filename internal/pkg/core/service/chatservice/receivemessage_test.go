package chatservice_test

import (
	"testing"
	"websocket_client/internal/pkg/core/adapter/kvadapter"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"
	"websocket_client/internal/pkg/core/adapter/wsconnadapter"
	"websocket_client/internal/pkg/core/adapter/wsstoreadapter"
	"websocket_client/internal/pkg/core/service/chatservice"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestReceiveMessage(t *testing.T) {
	// Define mocks
	mockCtrl := gomock.NewController(t)
	mockLogger := loggeradapter.NewMockAdapter(mockCtrl)
	mockRedisClient := kvadapter.NewMockRepoAdapter(mockCtrl)
	mockWsConn := wsconnadapter.NewMockAdapter(mockCtrl)
	mockWsStore := wsstoreadapter.NewMockAdapter(mockCtrl)

	type ReceiveMessageTestDataReq struct {
		userID string
	}

	type ReceiveMessageTestData struct {
		name       string
		req        ReceiveMessageTestDataReq
		mock       func()
		assertions func(isOnline bool, err error)
		isOnline   bool
	}

	tests := []ReceiveMessageTestData{
		{
			name: "Test user online and successful writing of json",
			req: ReceiveMessageTestDataReq{
				userID: "foo",
			},
			mock: func() {
				mockWsStore.EXPECT().GetConn(gomock.Any()).Return(mockWsConn).Times(1)
				mockWsConn.EXPECT().WriteJSON(gomock.Any()).Return(nil).Times(1)
				mockWsStore.EXPECT().Release(gomock.Any()).Times(1)
			},
			assertions: func(isOnline bool, err error) {
				assert.True(t, isOnline)
				assert.Nil(t, err)
			},
			isOnline: true,
		},
		{
			name: "Test user offline",
			req: ReceiveMessageTestDataReq{
				userID: "foo",
			},
			mock: func() {
				mockWsStore.EXPECT().GetConn(gomock.Any()).Return(nil).Times(1)
			},
			assertions: func(isOnline bool, err error) {
				assert.False(t, isOnline)
				assert.Nil(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			chatService := chatservice.NewChatService(chatservice.NewChatServiceReq{
				Logger: mockLogger, Redis: mockRedisClient, WsStore: mockWsStore,
			})

			test.mock()

			isOnline, err := chatService.ReceiveMessage(test.req.userID, test.req)

			test.assertions(isOnline, err)

		})
	}
}
