package chatservice

import (
	"errors"
	"testing"
	"websocket_client/internal/pkg/core/adapter/chatadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"

	"websocket_client/internal/pkg/platform/mysql/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestChatService_SearchUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		searchUserID   string
		mockDBResponse []models.Account
		mockDBError    error
		expectedResp   chatadapter.SearchUserResp
	}{
		{
			name:         "Successful User Search",
			searchUserID: "user1",
			mockDBResponse: []models.Account{
				{UserID: "user1"},
				{UserID: "user2"},
			},
			mockDBError: nil,
			expectedResp: chatadapter.SearchUserResp{
				UserIDs: []string{"user1", "user2"},
			},
		},
		{
			name:           "Failed User Search",
			searchUserID:   "user1",
			mockDBResponse: nil,
			mockDBError:    errors.New("database error"),
			expectedResp: chatadapter.SearchUserResp{
				UserIDs: []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := databaseadapter.NewMockRepoAdapter(ctrl)
			mockLogger := loggeradapter.NewMockAdapter(ctrl)

			// Expect DB call with given user ID and return mocked response/error.
			mockDB.EXPECT().SearchUser(databaseadapter.SearchUserRequest{
				UserID: tt.searchUserID,
			}).Return(tt.mockDBResponse, tt.mockDBError)

			// If there is an error, expect logger to log it.
			if tt.mockDBError != nil {
				mockLogger.EXPECT().NewError(tt.mockDBError.Error()).Times(1)
			}

			svc := NewChatService(NewChatServiceReq{
				DB:     mockDB,
				Logger: mockLogger,
			})

			req := chatadapter.SearchUserReq{SearchUserID: tt.searchUserID}

			resp := svc.SearchUser(req)

			assert.Equal(t, tt.expectedResp.UserIDs, resp.UserIDs)
		})
	}
}
