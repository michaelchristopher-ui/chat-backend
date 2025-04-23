package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"websocket_client/api/http/structs"
	"websocket_client/internal/pkg/core/adapter/accountadapter"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAccount(t *testing.T) {
	// Define mocks
	mockCtrl := gomock.NewController(t)
	mockAccountService := accountadapter.NewMockAdapter(mockCtrl)
	mockLogger := loggeradapter.NewMockAdapter(mockCtrl)

	// Define test data
	type RegisterAccountTestData struct {
		name       string
		wrongBody  bool
		reqBody    RegisterAccountReq
		mock       func()
		assertions func(res error, rec *httptest.ResponseRecorder)
	}

	// Populate test data
	tests := []RegisterAccountTestData{
		{
			name: "Test if correct body returns correct response provided the register function does not return an error",
			reqBody: RegisterAccountReq{
				UserID:   "a",
				Password: "a",
			},
			mock: func() {
				mockAccountService.EXPECT().Register(accountadapter.RegisterReq{
					UserID:   "a",
					Password: "a",
				}).Return(nil).Times(1)
				mockLogger.EXPECT().NewInfo(gomock.Any()).Times(1)

			},
			assertions: func(res error, rec *httptest.ResponseRecorder) {
				if assert.NoError(t, res) {
					assert.Equal(t, http.StatusOK, rec.Code)
					assert.Equal(t, "null\n", rec.Body.String())
				}
			},
		},
		{
			name: "Test if correct body returns error should the password be more than 72 characters",
			reqBody: RegisterAccountReq{
				UserID:   "a",
				Password: "CNKBVNRMXUEMSFSRLKJQBLYTQEWBGYCCKQGHZZMWDLNTCUMRYEXSCPPEHSGQHNLCQEBFGKWTJ",
			},
			mock: func() {
				mockLogger.EXPECT().NewInfo(gomock.Any()).Times(1)

			},
			assertions: func(res error, rec *httptest.ResponseRecorder) {
				if assert.NoError(t, res) {
					assert.Equal(t, http.StatusBadRequest, rec.Code)
					expected, _ := json.Marshal(structs.ErrorRet{
						Error: returnErrorPasswordLength,
					})
					assert.Equal(t, string(expected)+"\n", rec.Body.String())
				}
			},
		},
		{
			name: "Test if correct body returns error response provided the register function returns an error",
			reqBody: RegisterAccountReq{
				UserID:   "a",
				Password: "a",
			},
			mock: func() {
				mockAccountService.EXPECT().Register(accountadapter.RegisterReq{
					UserID:   "a",
					Password: "a",
				}).Return(errors.New("foo")).Times(1)
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
					//We don't handle the error string for this one, we just assume there will be an error.
					assert.True(t, bytes.Contains(rec.Body.Bytes(), []byte("\"error\"")))
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set up the body. Below, we handle one unique case first which is the invalid json before going to different requests controlled by the wrongBody boolean.
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte("{invalid_json}")))
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
					AccountService: mockAccountService,
					Logger:         mockLogger,
				},
			)

			// Execute mocks
			test.mock()

			// Run the tested function
			res := mockIntegrator.RegisterAccount(c)

			// Check the assertions
			test.assertions(res, rec)
		})
	}
}
