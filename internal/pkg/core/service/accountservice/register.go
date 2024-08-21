package accountservice

import (
	"websocket_client/internal/pkg/core/adapter/accountadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"

	"golang.org/x/crypto/bcrypt"
)

// Register encrypts the password with bcrypt with default cost before creating a new account entry in the db
func (a AccountService) Register(req accountadapter.RegisterReq) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return a.DB.SetAccount(databaseadapter.SetAccountReq{
		UserID:   req.UserID,
		Password: string(hashedPassword),
	})
}
