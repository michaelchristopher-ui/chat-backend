package accountservice

import (
	"fmt"
	"websocket_client/internal/pkg/core/adapter/accountadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
)

// Register encrypts the password with bcrypt with default cost (10) before creating a new account entry in the db
func (a AccountService) Register(req accountadapter.RegisterReq) error {
	hashedPassword, err := a.passwordGenerator.Generate(req.Password)
	if err != nil {
		a.logger.NewError(fmt.Sprintf(logErrRegisterFormat, "Error when hashing password", err.Error()))
		return err
	}

	err = a.db.SetAccount(databaseadapter.SetAccountReq{
		UserID:   req.UserID,
		Password: hashedPassword,
	})
	if err != nil {
		a.logger.NewError(fmt.Sprintf(logErrRegisterFormat, "Error when setting account", err.Error()))
		return err
	}
	a.logger.NewInfo(fmt.Sprintf(logInfoRegistered, req.UserID))
	return nil
}
