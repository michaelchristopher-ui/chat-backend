package main

import (
	"fmt"
	"sync"
	apihttp "websocket_client/api/http"
	"websocket_client/internal/common"
	"websocket_client/internal/conf"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"
	"websocket_client/internal/pkg/core/adapter/wsconnadapter"
	"websocket_client/internal/pkg/core/service/accountservice"
	"websocket_client/internal/pkg/core/service/chatservice"
	"websocket_client/internal/pkg/core/service/senderservice"
	"websocket_client/internal/pkg/core/service/wsprocessorservice"
	"websocket_client/internal/pkg/core/service/wsstoreservice"
	wsupgrader "websocket_client/internal/pkg/core/service/wsupgraderservice"
	"websocket_client/internal/pkg/platform/mysql"
	"websocket_client/internal/pkg/platform/redis"
	"websocket_client/internal/pkg/platform/zaplogger"
	"websocket_client/internal/transport"

	"github.com/gorilla/websocket"
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
				lgr.NewError(" Server panicked, err: %v", err)
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

	//Init Sub-Services
	senderService := senderservice.NewSenderService(senderservice.NewSenderServiceReq{
		DB:     db,
		Redis:  rds,
		Logger: lgr,
	})

	//Init Services
	chatService := chatservice.NewChatService(chatservice.NewChatServiceReq{
		DB:     db,
		Redis:  rds,
		Logger: lgr,
		Sender: senderService,
		WsStore: wsstoreservice.NewChatBackendService(wsstoreservice.NewWsStoreServiceReq{
			UserConnections: map[string]wsconnadapter.Adapter{},
			Lock:            map[string]*sync.Mutex{},
		}),
	})

	accountService := accountservice.NewAccountService(accountservice.NewAccountServiceReq{
		DB:     db,
		Logger: lgr,
	})

	wsstoreService := wsstoreservice.NewChatBackendService(wsstoreservice.NewWsStoreServiceReq{
		UserConnections: make(map[string]wsconnadapter.Adapter),
		Lock:            make(map[string]*sync.Mutex),
	})

	upgraderService := wsupgrader.NewWsUpgraderService(wsupgrader.NewWsUpgraderServiceReq{
		Upgrader: websocket.Upgrader{},
	})

	processorService := wsprocessorservice.NewWsProcessorService(wsprocessorservice.WsProcessorServiceReq{
		Logger:  lgr,
		WsStore: wsstoreService,
	})

	//Init server
	srv := transport.NewServer()

	//Init APIs
	apihttp.API(apihttp.APIReq{
		E:              srv.GetEcho(),
		ChatService:    chatService,
		AccountService: accountService,
		Logger:         lgr,
		WsProcessor:    processorService,
		Upgrader:       upgraderService,
	})

	//Start the server
	srv.StartServer()
}
