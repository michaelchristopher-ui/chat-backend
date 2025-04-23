package accountservice_test

import (
	"errors"
	"testing"
	"websocket_client/internal/pkg/core/adapter/accountadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"
	"websocket_client/internal/pkg/core/service/accountservice"
	"websocket_client/internal/pkg/platform/mysql/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAccountService_VerifyAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := databaseadapter.NewMockRepoAdapter(ctrl)
	mockLogger := loggeradapter.NewMockAdapter(ctrl)

	accountService := accountservice.NewAccountService(accountservice.NewAccountServiceReq{
		DB:     mockDB,
		Logger: mockLogger,
	})

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.DefaultCost)

	tests := []struct {
		name        string
		req         accountadapter.VerifyAuthReq
		setupMocks  func()
		expectedErr bool
	}{
		{
			name: "success",
			req: accountadapter.VerifyAuthReq{
				UserID:   "user1",
				Password: "correct_password",
			},
			setupMocks: func() {
				mockDB.EXPECT().
					GetAccount(databaseadapter.GetAccountReq{UserId: "user1"}).
					Return(models.Account{Password: string(hashedPassword)}, nil).Times(1)
			},
			expectedErr: false,
		},
		{
			name: "db error",
			req: accountadapter.VerifyAuthReq{
				UserID: "user2", Password: "any"},
			setupMocks: func() {
				mockDB.EXPECT().
					GetAccount(databaseadapter.GetAccountReq{UserId: "user2"}).
					Return(models.Account{}, errors.New("db error")).Times(1)
				mockLogger.EXPECT().NewError(gomock.Any(), gomock.Any()).Times(1)
			},
			expectedErr: true,
		},
		{
			name: "password mismatch",
			req: accountadapter.VerifyAuthReq{
				UserID: "user3", Password: "wrong_password"},
			setupMocks: func() {
				mockDB.EXPECT().
					GetAccount(databaseadapter.GetAccountReq{UserId: "user3"}).
					Return(models.Account{Password: string(hashedPassword)}, nil).Times(1)
				mockLogger.EXPECT().NewError(gomock.Any(), gomock.Any()).Times(1)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := accountService.VerifyAuth(tt.req)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
