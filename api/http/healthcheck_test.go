package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Define mocks
	mockCtrl := gomock.NewController(t)
	mockLogger := loggeradapter.NewMockAdapter(mockCtrl)

	// Define test data
	type HealthCheckTestData struct {
		name       string
		mock       func()
		assertions func(res error, rec *httptest.ResponseRecorder)
	}

	// Populate test data
	tests := []HealthCheckTestData{
		{
			name: "Test if HealthCheck returns correct response",
			mock: func() {
				mockLogger.EXPECT().NewInfo(gomock.Any(), gomock.Any()).Times(1)
			},
			assertions: func(res error, rec *httptest.ResponseRecorder) {
				if assert.NoError(t, res) {
					assert.Equal(t, http.StatusOK, rec.Code)
					assert.Equal(t, "null\n", rec.Body.String())
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a new HTTP request
			req := httptest.NewRequest(http.MethodGet, "/health", nil)

			// Create a new HTTP response recorder
			rec := httptest.NewRecorder()

			// Create a new context
			c := echo.New().NewContext(req, rec)

			// Initialize tested struct
			mockIntegrator := NewAPIIntegrator(
				NewAPIIntegratorReq{
					Logger: mockLogger,
				},
			)

			// Execute mocks
			test.mock()

			// Run the tested function
			res := mockIntegrator.HealthCheck(c)

			// Check the assertions
			test.assertions(res, rec)
		})
	}
}
