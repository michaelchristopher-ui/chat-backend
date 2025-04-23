package http

import (
	"websocket_client/internal/common"
	"websocket_client/internal/pkg/core/adapter/accountadapter"

	"github.com/labstack/echo/v4"
)

// AddContextID returns a function which generates FlowID and puts it within Context
func (integrator *APIIntegrator) AddContextID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		/*
			Generate FlowID for easy tracing
			The FlowID is a UUID. Do take note of the RPS to find out when we should reset or recycle the UUIDs.
		*/
		flowID := common.GenerateUUID()
		c.Set("flow_id", flowID)
		return next(c)
	}
}

// Auth is a rudimentary authorization function that checks whether the user exists in the database and the password is correct
func (integrator *APIIntegrator) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		defer func() {
			if err != nil {
				integrator.Logger.NewError(err.Error())
			}
		}()
		userId, password, err := common.SplitUserIDAndPasswordFromAuth(c.Request().Header.Get("Authorization"))
		if err != nil {
			c.Error(err)
			return err
		}

		err = integrator.AccountService.VerifyAuth(accountadapter.VerifyAuthReq{
			UserID:   userId,
			Password: password,
		})

		if err != nil {
			c.Error(err)
			return err
		}
		return next(c)
	}
}
