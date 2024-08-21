package accountservice

import (
	"websocket_client/internal/pkg/core/adapter/accountadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
)

// AccountService is a service that handles all account-related functionalities
type AccountService struct {
	DB databaseadapter.RepoAdapter
}

// NewAccountService is a constructor function for AccountService, which conforms to accountadapter.Adapter
func NewAccountService(DB databaseadapter.RepoAdapter) accountadapter.Adapter {
	return AccountService{
		DB: DB,
	}
}
