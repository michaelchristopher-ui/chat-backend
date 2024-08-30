package accountservice

import (
	"fmt"
	"websocket_client/internal/pkg/core/adapter/accountadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"

	"golang.org/x/crypto/bcrypt"
)

// VerifyAuth is a function to define whether or not the supplied password for the user is correct
func (a AccountService) VerifyAuth(req accountadapter.VerifyAuthReq) error {
	account, err := a.db.GetAccount(databaseadapter.GetAccountReq{
		UserId: req.UserID,
	})
	if err != nil {
		a.logger.NewError(fmt.Sprintf(logErrFormat, err.Error()))
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(req.Password))
	if err != nil {
		a.logger.NewError(fmt.Sprintf(logErrFormat, err.Error()))
		return err
	}

	return nil
}
