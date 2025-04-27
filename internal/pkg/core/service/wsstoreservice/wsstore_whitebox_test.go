package wsstoreservice

import (
	"sync"
	"testing"
	"time"
	"websocket_client/internal/pkg/core/adapter/lockeradapter"
	"websocket_client/internal/pkg/core/adapter/wsconnadapter"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteConn(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockLockMutex := lockeradapter.NewMockLocker(ctrl)
	mockFooMutex := lockeradapter.NewMockLocker(ctrl)

	mockConn := wsconnadapter.NewMockAdapter(ctrl)

	tests := []struct {
		name            string
		mock            func()
		userConnections map[string]wsconnadapter.Adapter
		connLock        map[string]sync.Locker
	}{
		{
			name: "Test DeleteConn locks lockMutex and user mutex, then deletes entries",
			mock: func() {
				mockLockMutex.EXPECT().Lock()
				mockFooMutex.EXPECT().Lock()
				mockLockMutex.EXPECT().Unlock()
			},
			userConnections: map[string]wsconnadapter.Adapter{
				"foo": mockConn,
			},
			connLock: map[string]sync.Locker{
				"foo": mockFooMutex,
			},
		},
		{
			name: "Test DeleteConn locks lockMutex but no user mutex exists",
			mock: func() {
				mockLockMutex.EXPECT().Lock()
				mockLockMutex.EXPECT().Unlock()
			},
			userConnections: map[string]wsconnadapter.Adapter{
				"foo": mockConn,
			},
			connLock: map[string]sync.Locker{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWsStore := NewWSStoreService(NewWsStoreServiceReq{
				ConnLock:        tt.connLock,
				LockMutex:       mockLockMutex,
				UserConnections: tt.userConnections,
			})

			tt.mock()
			mockWsStore.DeleteConn("foo")

			// Wait a short time to allow goroutines to complete
			time.Sleep(10 * time.Millisecond)

			// Assert the entries are deleted
			_, connExists := mockWsStore.(*WsStore).userConnections["foo"]
			_, lockExists := mockWsStore.(*WsStore).connLock["foo"]

			assert.False(t, connExists, "userConnections entry should be deleted")
			assert.False(t, lockExists, "connLock entry should be deleted")
		})
	}
}
