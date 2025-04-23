package accountservice_test

import (
	"errors"
	"testing"
	"websocket_client/internal/pkg/core/adapter/accountadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"
	"websocket_client/internal/pkg/core/adapter/passwordgeneratoradapter"
	"websocket_client/internal/pkg/core/service/accountservice"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPG := passwordgeneratoradapter.NewMockAdapter(ctrl)
	mockDB := databaseadapter.NewMockRepoAdapter(ctrl)
	mockLogger := loggeradapter.NewMockAdapter(ctrl)

	tests := []struct {
		name          string
		req           accountadapter.RegisterReq
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Successful Registration",
			req: accountadapter.RegisterReq{
				UserID:   "user1",
				Password: "password123",
			},
			mockSetup: func() {
				mockPG.EXPECT().Generate("password123").Return("hashed_password", nil).Times(1)
				mockDB.EXPECT().SetAccount(databaseadapter.SetAccountReq{
					UserID:   "user1",
					Password: "hashed_password",
				}).Return(nil).Times(1)
				mockLogger.EXPECT().NewInfo(gomock.Any()).Times(1)
			},
			expectedError: nil,
		},
		{
			name: "Password Generation Error",
			req: accountadapter.RegisterReq{
				UserID:   "user1",
				Password: "password123",
			},
			mockSetup: func() {
				mockPG.EXPECT().Generate("password123").Return("", errors.New("password generation failed")).Times(1)
				mockLogger.EXPECT().NewError(gomock.Any()).Times(1)
			},
			expectedError: errors.New("password generation failed"),
		},
		{
			name: "Database Set Account Error",
			req: accountadapter.RegisterReq{
				UserID:   "user1",
				Password: "password123",
			},
			mockSetup: func() {
				mockPG.EXPECT().Generate("password123").Return("hashed_password", nil).Times(1)
				mockDB.EXPECT().SetAccount(databaseadapter.SetAccountReq{
					UserID:   "user1",
					Password: "hashed_password",
				}).Return(errors.New("database error")).Times(1)
				mockLogger.EXPECT().NewError(gomock.Any()).Times(1)
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			accountService := accountservice.NewAccountService(accountservice.NewAccountServiceReq{
				PasswordGenerator: mockPG,
				DB:                mockDB,
				Logger:            mockLogger,
			})

			err := accountService.Register(tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
