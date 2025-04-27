package wsstoreservice_test

import (
	"sync"
	"testing"

	"websocket_client/internal/pkg/core/adapter/lockeradapter"
	"websocket_client/internal/pkg/core/adapter/wsconnadapter"
	"websocket_client/internal/pkg/core/service/wsstoreservice"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddConn(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockLockMutex := lockeradapter.NewMockLocker(ctrl)
	mockFooMutex := lockeradapter.NewMockLocker(ctrl)
	mockConnLock := map[string]sync.Locker{
		"foo": mockFooMutex,
	}

	mockConn := wsconnadapter.NewMockAdapter(ctrl)

	tests := []struct {
		name            string
		mock            func()
		userConnections map[string]wsconnadapter.Adapter
		expectError     bool
	}{
		{
			name: "Test AddConn creates and uses the mutexes properly before and after adding the new connection",
			mock: func() {
				mockLockMutex.EXPECT().Lock().Times(1)
				mockFooMutex.EXPECT().Lock().Times(1)
				mockFooMutex.EXPECT().Unlock().Times(1)
				mockLockMutex.EXPECT().Unlock().Times(1)
			},
			userConnections: map[string]wsconnadapter.Adapter{},
		},
		{
			name: "Test AddConn creates and uses the mutexes properly before and after adding the new connection, but connection exist",
			mock: func() {
				mockLockMutex.EXPECT().Lock().Times(1)
				mockFooMutex.EXPECT().Lock().Times(1)
				mockFooMutex.EXPECT().Unlock().Times(1)
				mockLockMutex.EXPECT().Unlock().Times(1)
			},
			userConnections: map[string]wsconnadapter.Adapter{
				"foo": mockConn,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWsStore := wsstoreservice.NewWSStoreService(wsstoreservice.NewWsStoreServiceReq{
				ConnLock:        mockConnLock,
				LockMutex:       mockLockMutex,
				UserConnections: tt.userConnections,
			})

			tt.mock()
			err := mockWsStore.AddConn(mockConn, "foo")
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockLockMutex := lockeradapter.NewMockLocker(ctrl)
	mockFooMutex := lockeradapter.NewMockLocker(ctrl)

	mockConn := wsconnadapter.NewMockAdapter(ctrl)

	tests := []struct {
		name            string
		mock            func()
		userConnections map[string]wsconnadapter.Adapter
		connLock        map[string]sync.Locker
		expectedConn    wsconnadapter.Adapter
	}{
		{
			name: "Test GetConn locks mutex and returns existing connection",
			mock: func() {
				mockLockMutex.EXPECT().Lock().Times(1)
				mockFooMutex.EXPECT().Lock().Times(1)
				mockLockMutex.EXPECT().Unlock().Times(1)
			},
			userConnections: map[string]wsconnadapter.Adapter{
				"foo": mockConn,
			},
			connLock: map[string]sync.Locker{
				"foo": mockFooMutex,
			},
			expectedConn: mockConn,
		},
		{
			name: "Test GetConn returns nil if no mutex for userID",
			mock: func() {
				mockLockMutex.EXPECT().Lock().Times(1)
				mockLockMutex.EXPECT().Unlock().Times(1)
			},
			userConnections: map[string]wsconnadapter.Adapter{},
			connLock:        map[string]sync.Locker{},
			expectedConn:    nil,
		},
		{
			name: "Test GetConn returns nil if no connection for userID",
			mock: func() {
				mockLockMutex.EXPECT().Lock()
				mockFooMutex.EXPECT().Lock()
				mockLockMutex.EXPECT().Unlock()
			},
			connLock: map[string]sync.Locker{
				"foo": mockFooMutex,
			},
			userConnections: map[string]wsconnadapter.Adapter{},
			expectedConn:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWsStore := wsstoreservice.NewWSStoreService(wsstoreservice.NewWsStoreServiceReq{
				ConnLock:        tt.connLock,
				LockMutex:       mockLockMutex,
				UserConnections: tt.userConnections,
			})

			tt.mock()
			conn := mockWsStore.GetConn("foo")
			assert.Equal(t, tt.expectedConn, conn)
		})
	}
}

func TestRelease(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockLockMutex := lockeradapter.NewMockLocker(ctrl)
	mockFooMutex := lockeradapter.NewMockLocker(ctrl)
	mockLock := map[string]sync.Locker{
		"foo": mockFooMutex,
	}

	tests := []struct {
		name string
		mock func()
	}{
		{
			name: "Test Release locks lockMutex and unlocks user mutex",
			mock: func() {
				mockLockMutex.EXPECT().Lock()
				mockFooMutex.EXPECT().Unlock()
				mockLockMutex.EXPECT().Unlock()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWsStore := wsstoreservice.NewWSStoreService(wsstoreservice.NewWsStoreServiceReq{
				ConnLock:  mockLock,
				LockMutex: mockLockMutex,
			})

			tt.mock()
			mockWsStore.Release("foo")
		})
	}
}
