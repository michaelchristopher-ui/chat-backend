package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"websocket_client/internal/pkg/core/adapter/chatadapter"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"
	"websocket_client/internal/pkg/core/adapter/wsprocessoradapter"
	"websocket_client/internal/pkg/core/adapter/wsupgraderadapter"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCallChatServiceFunc(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockChatService := chatadapter.NewMockAdapter(mockCtrl)
	mockLogger := loggeradapter.NewMockAdapter(mockCtrl)
	integrator := &APIIntegrator{
		ChatService: mockChatService,
		Logger:      mockLogger,
	}

	tests := []struct {
		name        string
		dataByte    []byte
		reqType     string
		expectError bool
		setupMock   func()
	}{
		{
			name:        "Successful Message Sending",
			dataByte:    json.RawMessage(`{"message":"Hello","from_user_id":"user1","to_user_id":"user2"}`),
			reqType:     incomingMessageTypeMessage,
			expectError: false,
			setupMock: func() {
				// Set up expectation for successful message sending
				mockChatService.EXPECT().SendMessage(gomock.Any()).Times(1)
				mockChatService.EXPECT().ReceiveMessage("user1", gomock.Any()).Return(true, nil).Times(1)
			},
		},
		{
			name:        "Successful Add Friend",
			dataByte:    json.RawMessage(`{"friend_id":"user1"}`),
			reqType:     incomingMessageTypeAddFriend,
			expectError: false,
			setupMock: func() {
				// Set up expectation for successful message sending
				mockChatService.EXPECT().AddFriend(gomock.Any()).Times(1)
				mockChatService.EXPECT().ReceiveMessage("user1", gomock.Any()).Return(true, nil).Times(1)
			},
		},
		{
			name:        "Successful Get Chat History",
			dataByte:    json.RawMessage(`{"from_user_id":"user1","to_user_id":"user2","offset":1,"limit":1,"timestamp_after":""}`),
			reqType:     incomingMessageTypeGetChatHistory,
			expectError: false,
			setupMock: func() {
				// Set up expectation for successful message sending
				mockChatService.EXPECT().GetChatHistory(gomock.Any()).Times(1)
				mockChatService.EXPECT().ReceiveMessage("user1", gomock.Any()).Return(true, nil).Times(1)
			},
		},
		{
			name:        "Successful Remove Friend",
			dataByte:    json.RawMessage(`{"friend_id":"user1"}`),
			reqType:     incomingMessageTypeRemoveFriend,
			expectError: false,
			setupMock: func() {
				// Set up expectation for successful message sending
				mockChatService.EXPECT().RemoveFriend(gomock.Any()).Times(1)
				mockChatService.EXPECT().ReceiveMessage("user1", gomock.Any()).Return(true, nil).Times(1)
			},
		},
		{
			name:        "Successful Search User",
			dataByte:    json.RawMessage(`{"search_user_id":"user1"}`),
			reqType:     incomingMessageTypeSearchUser,
			expectError: false,
			setupMock: func() {
				// Set up expectation for successful message sending
				mockChatService.EXPECT().SearchUser(gomock.Any()).Times(1)
				mockChatService.EXPECT().ReceiveMessage("user1", gomock.Any()).Return(true, nil).Times(1)
			},
		},
		{
			name:        "Invalid Request",
			dataByte:    []byte(`{}`),
			reqType:     "deliberately incorrect request type",
			expectError: true,
			setupMock: func() {
				mockLogger.EXPECT().NewError(gomock.Any(), gomock.Any()).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.setupMock != nil {
				tt.setupMock()
			}

			err := integrator.CallChatServiceFunc(tt.dataByte, "user1", tt.reqType)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandleWebsocket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup mocks for upgrader and logger based on test case.
	mockUpgrader := wsupgraderadapter.NewMockAdapter(ctrl)

	mockLogger := loggeradapter.NewMockAdapter(ctrl)
	mockWsProcessor := wsprocessoradapter.NewMockAdapter(ctrl)

	tests := []struct {
		name        string
		userID      string
		mock        func()
		expectError bool
	}{
		{
			name:   "Successful WebSocket Upgrade",
			userID: "user1",
			mock: func() {
				mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				mockWsProcessor.EXPECT().ProcessWebsocketAfterUpgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			integrator := NewAPIIntegrator(NewAPIIntegratorReq{
				Logger:      mockLogger,
				Upgrader:    mockUpgrader,
				WsProcessor: mockWsProcessor,
			})

			tt.mock()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/ws", nil)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			if tt.userID != "" {
				c.Set("user_id", tt.userID)
			}

			err := integrator.HandleWebsocket(c)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
