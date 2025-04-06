package main

import (
	"fmt"
	apihttp "websocket_client/api/http"
	"websocket_client/internal/common"
	"websocket_client/internal/conf"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"
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
	startServer()
}

func startServer() {
	var lgr loggeradapter.Adapter

	// Panic catcher before restarting server if error
	defer func() {
		if err := recover(); err != nil {
			if lgr != nil {
				lgr.NewError(fmt.Sprintf("[StartServer] Server panicked, err: %v", err))
			}
			startServer()
		}
	}()
	//Init custom zap logger
	lgr, err := zaplogger.NewLogger()
	if err != nil {
		panic(fmt.Sprintf("error setting up logger, err: %s", err.Error()))
	}

	//Init Configs
	err = conf.Init(*common.CfgPath)
	if err != nil {
		panic(fmt.Sprintf("error parsing config, err: %s", err.Error()))
	}

	//Init Components of Services
	db, err := mysql.NewDatabase()
	if err != nil {
		panic(fmt.Sprintf("error setting up database, err: %s", err.Error()))
	}

	rds, err := redis.NewRedis()
	if err != nil {
		panic(fmt.Sprintf("error setting up redis, err: %s", err.Error()))
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

	//Start the server
	srv.StartServer()
}
