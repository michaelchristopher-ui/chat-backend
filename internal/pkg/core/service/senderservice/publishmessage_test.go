package senderservice_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"websocket_client/internal/pkg/core/adapter/kvadapter"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"
	"websocket_client/internal/pkg/core/adapter/senderadapter"
	"websocket_client/internal/pkg/core/service/senderservice"

	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPublishMessage(t *testing.T) {
	// Define mocks
	mockCtrl := gomock.NewController(t)
	mockLogger := loggeradapter.NewMockAdapter(mockCtrl)
	mockRedisClient := kvadapter.NewMockRepoAdapter(mockCtrl)

	// Declare server
	var server *httptest.Server

	// Define test data
	type PublishMessageTestData struct {
		name           string
		req            senderadapter.SendMessageReq
		mock           func()
		assertions     func(isOnline bool, err error)
		serverResponse func(w http.ResponseWriter)
	}

	tests := []PublishMessageTestData{
		{
			name: "Test successful message publication when user is online",
			req: senderadapter.SendMessageReq{
				Message:  "Hello",
				ToUserID: "user123",
				Type:     1,
			},
			mock: func() {
				mockRedisClient.EXPECT().GetValue("user123").Return(server.URL[7:], nil).Times(1) // Simulate user being online
				mockLogger.EXPECT().NewInfo(gomock.Any()).Times(1)
			},
			assertions: func(isOnline bool, err error) {
				assert.True(t, isOnline)
				assert.NoError(t, err)
			},
			serverResponse: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte(`{"is_online":true}`)); err != nil {
					assert.FailNow(t, err.Error())
				}
			},
		},
		{
			name: "Test error returned by server",
			req: senderadapter.SendMessageReq{
				Message:  "Hello",
				ToUserID: "user123",
				Type:     1,
			},
			mock: func() {
				mockRedisClient.EXPECT().GetValue("user123").Return(server.URL[7:], nil).Times(1) // Simulate user being online
				mockLogger.EXPECT().NewInfo(gomock.Any()).Times(1)
			},
			assertions: func(isOnline bool, err error) {
				assert.Error(t, err)
			},
			serverResponse: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte(`{"error": "foo"}`)); err != nil {
					assert.FailNow(t, err.Error())
				}
			},
		},
		{
			name: "Test not convertible error returned by server",
			req: senderadapter.SendMessageReq{
				Message:  "Hello",
				ToUserID: "user123",
				Type:     1,
			},
			mock: func() {
				mockRedisClient.EXPECT().GetValue("user123").Return(server.URL[7:], nil).Times(1) // Simulate user being online
				mockLogger.EXPECT().NewInfo(gomock.Any()).Times(1)
			},
			assertions: func(isOnline bool, err error) {
				assert.EqualError(t, err, senderservice.ErrorExistsButNotConvertible.Error())
			},
			serverResponse: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte(`{"error": true}`)); err != nil {
					assert.FailNow(t, err.Error())
				}
			},
		},
		{
			name: "Test successful message publication when user is online",
			req: senderadapter.SendMessageReq{
				Message:  "Hello",
				ToUserID: "user123",
				Type:     1,
			},
			mock: func() {
				mockRedisClient.EXPECT().GetValue("user123").Return(server.URL[7:], nil).Times(1) // Simulate user being online
				mockLogger.EXPECT().NewInfo(gomock.Any()).Times(1)
			},
			assertions: func(isOnline bool, err error) {
				assert.True(t, isOnline)
				assert.NoError(t, err)
			},
			serverResponse: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte(`{"is_online":true}`)); err != nil {
					assert.FailNow(t, err.Error())
				}
			},
		},
		{
			name: "Test error when doing request",
			req: senderadapter.SendMessageReq{
				Message:  "Hello",
				ToUserID: "user123",
				Type:     1,
			},
			mock: func() {
				mockRedisClient.EXPECT().GetValue("user123").Return("foo", nil).Times(1)
				mockLogger.EXPECT().NewError(gomock.Any()).Times(1)
			},
			assertions: func(isOnline bool, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "Test successful message publication when user is online, but with a return that is undecodeable",
			req: senderadapter.SendMessageReq{
				Message:  "Hello",
				ToUserID: "user123",
				Type:     1,
			},
			mock: func() {
				mockRedisClient.EXPECT().GetValue("user123").Return(server.URL[7:], nil).Times(1) // Simulate user being online
				mockLogger.EXPECT().NewError(gomock.Any()).Times(1)
			},
			assertions: func(isOnline bool, err error) {
				assert.Error(t, err)
			},
			serverResponse: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte(`f`)); err != nil {
					assert.FailNow(t, err.Error())
				}
			},
		},
		{
			name: "Test IP retrieval returns an error",
			req: senderadapter.SendMessageReq{
				Message:  "Hello",
				ToUserID: "user123",
				Type:     1,
			},
			mock: func() {
				mockRedisClient.EXPECT().GetValue("user123").Return("", errors.New("foo")).Times(1) // Simulate user being online
				mockLogger.EXPECT().NewError(gomock.Any()).Times(1)
			},
			assertions: func(isOnline bool, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "Test IP retrieval returns empty ip",
			req: senderadapter.SendMessageReq{
				Message:  "Hello",
				ToUserID: "user123",
				Type:     1,
			},
			mock: func() {
				mockRedisClient.EXPECT().GetValue("user123").Return("", nil).Times(1) // Simulate user being online
				mockLogger.EXPECT().NewError(gomock.Any()).Times(1)
			},
			assertions: func(isOnline bool, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "Test if user offline returns an appropriate response.",
			req: senderadapter.SendMessageReq{
				Message:  "Hello!",
				Type:     2,
				ToUserID: "user123",
			},
			mock: func() {
				mockRedisClient.EXPECT().GetValue("user123").Return("", redis.Nil).Times(1) // Simulate user offline scenario.
				mockLogger.EXPECT().NewError(gomock.Any()).Times(1)
			},
			assertions: func(isOnline bool, err error) {
				assert.False(t, isOnline)
				assert.NoError(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/receive" {
					t.Errorf("Expected to request '/receive', got: %s", r.URL.Path)
				}

				test.serverResponse(w)
			}))

			defer server.Close()

			senderService := senderservice.NewSenderService(senderservice.NewSenderServiceReq{Logger: mockLogger, Redis: mockRedisClient})

			test.mock()

			isOnline, err := senderService.PublishMessage(test.req)

			test.assertions(isOnline, err)

		})
	}
}
