package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"websocket_client/internal/pkg/core/adapter/accountadapter"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSetFlowID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the function under test
	dummyNextFunc := func(c echo.Context) error {
		return nil
	}

	integrator := NewAPIIntegrator(NewAPIIntegratorReq{})

	handlerFunc := integrator.AddContextID(dummyNextFunc)

	err := handlerFunc(c)
	assert.NoError(t, err)

	val := c.Get("flow_id")
	assert.NotNil(t, val)

	strVal, ok := val.(string)
	assert.True(t, ok)

	// Optionally validate if strVal is a valid UUID (using github.com/google/uuid or similar)
	if _, err = uuid.Parse(strVal); err != nil {
		t.Errorf("Expected valid UUID for flow_id but got: %v", strVal)
	}
}

func TestAuthenticationHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAccountService := accountadapter.NewMockAdapter(ctrl)
	mockLogger := loggeradapter.NewMockAdapter(ctrl)

	tests := []struct {
		name               string
		authHeader         string
		mock               func()
		expectError        bool
		expectedStatusCode int
		expectLogCall      bool
	}{
		{
			name:       "Successful Authentication",
			authHeader: "Basic dXNlcjE6cGFzc3dvcmQx", // base64 user1:password1,
			mock: func() {
				mockAccountService.EXPECT().VerifyAuth(gomock.Any()).Return(nil).Times(1)
			},
			expectError:        false,
			expectedStatusCode: http.StatusOK,
			expectLogCall:      false,
		},
		{
			name:       "Invalid Authorization Header",
			authHeader: "InvalidHeader",
			mock: func() {
				mockLogger.EXPECT().NewError(gomock.Any()).Times(1)
			},
			expectError:        true,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:       "Failed Verification",
			authHeader: "Basic dXNlcjE6cGFzc3dvcmQx",
			mock: func() {
				mockAccountService.EXPECT().VerifyAuth(gomock.Any()).Return(errors.New("foo")).Times(1)
				mockLogger.EXPECT().NewError(gomock.Any()).Times(1)
			},
			expectError:        true,
			expectedStatusCode: http.StatusInternalServerError,
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks based on test case:
			integrator := NewAPIIntegrator(
				NewAPIIntegratorReq{
					AccountService: mockAccountService,
					Logger:         mockLogger,
				},
			)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", tt.authHeader)

			rec := httptest.NewRecorder()

			tt.mock()

			c := e.NewContext(req, rec)

			dummyNextFunc := func(c echo.Context) error {
				return nil
			}

			handlerFunc := integrator.Auth(dummyNextFunc)
			err := handlerFunc(c)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatusCode, rec.Code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatusCode, rec.Code)
			}
		})

	}

}
