package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"websocket_client/api/http/structs"
	"websocket_client/internal/pkg/core/adapter/chatadapter"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestReceiveMessage(t *testing.T) {
	// Define mocks
	mockCtrl := gomock.NewController(t)
	mockChatService := chatadapter.NewMockAdapter(mockCtrl)
	mockLogger := loggeradapter.NewMockAdapter(mockCtrl)

	// Define test data
	type ReceiveMessageTestData struct {
		name       string
		wrongBody  bool
		reqBody    ReceiveMessageReq
		mock       func()
		assertions func(res error, rec *httptest.ResponseRecorder)
	}

	// Populate test data
	tests := []ReceiveMessageTestData{
		{
			name: "Test if correct body returns correct response provided the user is online",
			reqBody: ReceiveMessageReq{
				Message:    "Hello",
				FromUserID: "user1",
				Type:       1,
				ToUserID:   "user2",
				Timestamp:  "2025-04-02T06:53:00Z",
			},
			mock: func() {
				mockChatService.EXPECT().ReceiveMessage(chatadapter.ReceiveMessageReq{
					Message:    "Hello",
					FromUserID: "user1",
					Type:       1,
					ToUserID:   "user2",
					Timestamp:  "2025-04-02T06:53:00Z",
				}).Return(true, nil).Times(1)
				mockLogger.EXPECT().NewInfo(gomock.Any()).Times(1)
			},
			assertions: func(res error, rec *httptest.ResponseRecorder) {
				if assert.NoError(t, res) {
					assert.Equal(t, http.StatusOK, rec.Code)
					expected, _ := json.Marshal(ReceiveMessageRes{
						IsOnline: true,
					})
					assert.Equal(t, string(expected)+"\n", rec.Body.String())
				}
			},
		},
		{
			name: "Test if correct body returns correct response provided the user is not online",
			reqBody: ReceiveMessageReq{
				Message:    "Hello",
				FromUserID: "user1",
				Type:       1,
				ToUserID:   "user2",
				Timestamp:  "2025-04-02T06:53:00Z",
			},
			mock: func() {
				mockChatService.EXPECT().ReceiveMessage(chatadapter.ReceiveMessageReq{
					Message:    "Hello",
					FromUserID: "user1",
					Type:       1,
					ToUserID:   "user2",
					Timestamp:  "2025-04-02T06:53:00Z",
				}).Return(false, nil).Times(1)
				mockLogger.EXPECT().NewInfo(gomock.Any()).Times(1)
			},
			assertions: func(res error, rec *httptest.ResponseRecorder) {
				if assert.NoError(t, res) {
					assert.Equal(t, http.StatusOK, rec.Code)
					expected, _ := json.Marshal(ReceiveMessageRes{
						IsOnline: false,
					})
					assert.Equal(t, string(expected)+"\n", rec.Body.String())
				}
			},
		},
		{
			name: "Test if correct body returns error response provided the chat service returns an error",
			reqBody: ReceiveMessageReq{
				Message:    "Hello",
				FromUserID: "user1",
				Type:       1,
				ToUserID:   "user2",
				Timestamp:  "2025-04-02T06:53:00Z",
			},
			mock: func() {
				mockChatService.EXPECT().ReceiveMessage(chatadapter.ReceiveMessageReq{
					Message:    "Hello",
					FromUserID: "user1",
					Type:       1,
					ToUserID:   "user2",
					Timestamp:  "2025-04-02T06:53:00Z",
				}).Return(false, errors.New("foo")).Times(1)
				mockLogger.EXPECT().NewError(gomock.Any()).Times(1)
			},
			assertions: func(res error, rec *httptest.ResponseRecorder) {
				if assert.NoError(t, res) {
					assert.Equal(t, http.StatusInternalServerError, rec.Code)
					expected, _ := json.Marshal(structs.ErrorRet{
						Error: "foo",
					})
					assert.Equal(t, string(expected)+"\n", rec.Body.String())
				}
			},
		},
		{
			name:      "Test if wrong body returns error",
			wrongBody: true,
			mock: func() {
				mockLogger.EXPECT().NewError(gomock.Any()).Times(1)
			},
			assertions: func(res error, rec *httptest.ResponseRecorder) {
				if assert.NoError(t, res) {
					assert.Equal(t, http.StatusInternalServerError, rec.Code)
					// We don't handle the error string for this one, we just assume there will be an error.
					assert.True(t, bytes.Contains(rec.Body.Bytes(), []byte("\"error\"")))
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set up the body. Below, we handle one unique case first which is the invalid json before going to different requests controlled by the wrongBody boolean.
			var req *http.Request = httptest.NewRequest(http.MethodPost, "/receive_message", bytes.NewReader([]byte("{invalid_json}")))
			if !test.wrongBody {
				jsonData, err := json.Marshal(test.reqBody)
				assert.NoError(t, err)
				req = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
			}

			// Create a new HTTP request

			rec := httptest.NewRecorder()

			// Create a new context
			c := echo.New().NewContext(req, rec)

			// Initialize tested struct
			mockIntegrator := NewAPIIntegrator(
				NewAPIIntegratorReq{
					ChatService: mockChatService,
					Logger:      mockLogger,
				},
			)

			// Execute mocks
			test.mock()

			// Run the tested function
			res := mockIntegrator.ReceiveMessage(c)

			// Check the assertions
			test.assertions(res, rec)
		})
	}
}
