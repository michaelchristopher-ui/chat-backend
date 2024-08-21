package accountservice

import (
	"websocket_client/internal/pkg/core/adapter/accountadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"

	"golang.org/x/crypto/bcrypt"
)

// VerifyAuth is a function to define whether or not the supplied password for the user is correct
func (a AccountService) VerifyAuth(req accountadapter.VerifyAuthReq) (err error) {
	account, err := a.DB.GetAccount(databaseadapter.GetAccountReq{
		UserId: req.UserID,
	})
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(req.Password))
	return
}
