package chatservice_test

import (
	"errors"
	"testing"
	"websocket_client/internal/pkg/core/adapter/chatadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"
	"websocket_client/internal/pkg/core/service/chatservice"

	"github.com/golang/mock/gomock"
)

func TestRemoveFriend(t *testing.T) {
	// Set up mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := databaseadapter.NewMockRepoAdapter(ctrl)
	mockLogger := loggeradapter.NewMockAdapter(ctrl)
	// Define test data
	type RemoveFriendTestData struct {
		name            string
		removeFriendReq chatadapter.RemoveFriendReq
		mock            func()
	}

	// Populate test data
	tests := []RemoveFriendTestData{
		{
			name: "Successful Removal",
			removeFriendReq: chatadapter.RemoveFriendReq{
				UserID:   "user1",
				FriendID: "friend1",
			},
			mock: func() {
				mockDB.EXPECT().RemoveFriend(databaseadapter.RemoveFriendReq{
					UserID:   "user1",
					FriendID: "friend1",
				}).Return(nil)
				mockLogger.EXPECT().NewInfo(gomock.Any(), gomock.Any()).Times(1)
			},
		},
		{
			name: "Failed Removal",
			removeFriendReq: chatadapter.RemoveFriendReq{
				UserID:   "user1",
				FriendID: "friend1",
			},
			mock: func() {
				mockDB.EXPECT().RemoveFriend(databaseadapter.RemoveFriendReq{
					UserID:   "user1",
					FriendID: "friend1",
				}).Return(errors.New("failed to remove friend"))
				mockLogger.EXPECT().NewError(gomock.Any(), gomock.Any()).Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// Initialize ChatService with mocks
			chatService := chatservice.NewChatService(chatservice.NewChatServiceReq{
				DB:     mockDB,
				Logger: mockLogger,
			})

			// Execute mocks
			test.mock()

			// Run the tested function
			chatService.RemoveFriend(test.removeFriendReq)
		})
	}
}
