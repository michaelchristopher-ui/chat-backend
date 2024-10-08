package main

import (
	"fmt"
	apihttp "websocket_client/api/http"
	"websocket_client/internal/common"
	"websocket_client/internal/conf"
	"websocket_client/internal/pkg/core/service/accountservice"
	"websocket_client/internal/pkg/core/service/chatservice"
	"websocket_client/internal/pkg/platform/mysql"
	"websocket_client/internal/pkg/platform/redis"
	"websocket_client/internal/pkg/platform/zaplogger"
	"websocket_client/internal/transport"
)

func main() {
	//Environment Variables
	common.SetEnvVars()

	err := conf.Init(*common.CfgPath)
	if err != nil {
		panic(fmt.Errorf("error parsing config. %w", err))
	}

	//Init Components of Services
	db, err := mysql.NewDatabase()
	if err != nil {
		panic(fmt.Errorf("error setting up database. %w", err))
	}

	rds, err := redis.NewRedis()
	if err != nil {
		panic(fmt.Errorf("error setting up redis. %w", err))
	}

	lgr, err := zaplogger.NewLogger()
	if err != nil {
		panic(fmt.Errorf("error setting up logger. %w", err))
	}

	//Init Services
	chatService := chatservice.NewChatService(chatservice.NewChatServiceReq{
		DB:     db,
		Redis:  rds,
		Logger: lgr,
	})

	accountService := accountservice.NewAccountService(accountservice.NewAccountServiceReq{
		DB:     db,
		Logger: lgr,
	})

	//Init server
	srv := transport.NewServer()

	//Init APIs
	apihttp.API(apihttp.APIReq{
		E:              srv.GetEcho(),
		ChatService:    chatService,
		AccountService: accountService,
		Logger:         lgr,
	})

	srv.StartServer()
}
