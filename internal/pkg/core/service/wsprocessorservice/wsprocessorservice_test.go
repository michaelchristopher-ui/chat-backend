package wsprocessorservice

import (
	"errors"
	"testing"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"
	"websocket_client/internal/pkg/core/adapter/wsconnadapter"
	"websocket_client/internal/pkg/core/adapter/wsstoreadapter"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestProcessWebsocketAfterUpgrade(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWSConn := wsconnadapter.NewMockAdapter(ctrl)

	mockWsStore := wsstoreadapter.NewMockAdapter(ctrl)
	mockLogger := loggeradapter.NewMockAdapter(ctrl)

	tests := []struct {
		name                 string
		mock                 func()
		expectedError        error
		isConnectionNil      bool
		IsCallServiceFuncNil bool
	}{
		{
			name:            "Connection is nil",
			isConnectionNil: true,
			expectedError:   errors.New(errNilConnStr),
		},
		{
			name:                 "Call Service Func is nil",
			IsCallServiceFuncNil: true,
			expectedError:        errors.New(errNilCallServiceStr),
		},
		{
			name: "Add conn failure",
			mock: func() {
				mockWsStore.EXPECT().AddConn(gomock.Any(), gomock.Any()).Return(errors.New("string")).Times(1)
				mockWsStore.EXPECT().DeleteConn(gomock.Any()).Times(1)
				mockWSConn.EXPECT().Close().Times(1)
			},
		},
		{
			name: "Successful Message Reading once then close",
			mock: func() {
				mockWsStore.EXPECT().AddConn(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				mockWSConn.EXPECT().ReadMessage().Return(
					websocket.TextMessage,
					[]byte(`{"req_type":"test","data":123}`),
					nil,
				).Times(1)
				mockLogger.EXPECT().NewInfo(gomock.Any(), gomock.Any())
				mockWSConn.EXPECT().ReadMessage().Return(
					websocket.CloseMessage,
					[]byte(``),
					nil,
				).Times(1)
				mockWsStore.EXPECT().DeleteConn(gomock.Any()).Times(1)
				mockWSConn.EXPECT().Close().Times(1)
			},
		},
		{
			name: "Fail reading message with error",
			mock: func() {
				mockWsStore.EXPECT().AddConn(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				mockWSConn.EXPECT().ReadMessage().Return(
					websocket.TextMessage,
					[]byte(``),
					errors.New("foo"),
				).Times(1)
				mockLogger.EXPECT().NewError(gomock.Any(), gomock.Any()).Times(1)
				mockWsStore.EXPECT().DeleteConn(gomock.Any()).Times(1)
				mockWSConn.EXPECT().Close().Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WsProcessor := NewWsProcessorService(WsProcessorServiceReq{
				Logger:  mockLogger,
				WsStore: mockWsStore,
			})

			if tt.mock != nil {
				tt.mock()
			}

			var wsConnParam wsconnadapter.Adapter
			if !tt.isConnectionNil {
				wsConnParam = mockWSConn
			}

			var dummyFuncParmm func(dataByte []byte, userID, reqType string) error = nil
			if !tt.IsCallServiceFuncNil {
				dummyFuncParmm = func(dataByte []byte, userID, reqType string) error {
					return nil
				}
			}

			err := WsProcessor.ProcessWebsocketAfterUpgrade(wsConnParam, "user1", dummyFuncParmm)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
