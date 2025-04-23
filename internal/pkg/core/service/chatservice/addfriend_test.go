package chatservice_test

import (
	"errors"
	"testing"
	"websocket_client/internal/pkg/core/adapter/chatadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
	"websocket_client/internal/pkg/core/adapter/kvadapter"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"
	"websocket_client/internal/pkg/core/service/chatservice"

	"github.com/golang/mock/gomock"
)

func TestAddFriend(t *testing.T) {
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

	tests := []struct {
		name         string
		addFriendReq chatadapter.AddFriendReq
		mock         func()
	}{
		{
			name: "Successful Friend Addition",
			addFriendReq: chatadapter.AddFriendReq{
				UserID:   "user123",
				FriendID: "friend456",
			},
			mock: func() {
				// Expect AddFriend to be called with correct parameters
				mockDB.EXPECT().AddFriend(databaseadapter.AddFriendReq{
					UserID:   "user123",
					FriendID: "friend456",
				}).Return(nil) // No error returned on success
			},
		},
		{
			name: "Error Adding Friend",
			addFriendReq: chatadapter.AddFriendReq{
				UserID:   "user123",
				FriendID: "friend456",
			},
			mock: func() {
				expectedErr := errors.New("database error")
				mockDB.EXPECT().AddFriend(databaseadapter.AddFriendReq{
					UserID:   "user123",
					FriendID: "friend456",
				}).Return(expectedErr) // Return an expected errror from db call
				mockLogger.EXPECT().NewError(gomock.Any(), gomock.Any())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mock() // Call mocked functions

			c.AddFriend(tt.addFriendReq) // Call AddFriends method

		})
	}
}
