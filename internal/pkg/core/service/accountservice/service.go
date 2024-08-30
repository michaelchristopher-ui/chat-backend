package accountservice

import (
	"websocket_client/internal/common"
	"websocket_client/internal/pkg/core/adapter/accountadapter"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"
)

// AccountService is a service that handles all account-related functionalities
type AccountService struct {
	db     databaseadapter.RepoAdapter
	logger loggeradapter.Adapter
}

// NewAccountService is a constructor function for AccountService, which conforms to accountadapter.Adapter
func NewAccountService(req NewAccountServiceReq) accountadapter.Adapter {
	if err := common.CheckNilFields(req); err != nil {
		panic(err.Error())
	}

	return AccountService{
		db:     req.DB,
		logger: req.Logger,
	}
}

// NewAccountServiceReq is a parameter struct for the NewAccountService function
type NewAccountServiceReq struct {
	DB     databaseadapter.RepoAdapter
	Logger loggeradapter.Adapter
}
