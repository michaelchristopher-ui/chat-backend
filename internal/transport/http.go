package transport

import (
	"net/http"
	"os"
	"time"

	"websocket_client/internal/common"
	config "websocket_client/internal/conf"

	"github.com/labstack/echo"
)

type server struct {
	e            *echo.Echo
	ipport       string
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewServer() server {
	e := echo.New()

	cfg := config.GetConfig()

	return server{
		e:            e,
		ipport:       *common.IPPort,
		readTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		writeTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}
}

func (h server) GetEcho() *echo.Echo {
	return h.e
}

func (h server) StartServer() {
	s := &http.Server{
		Addr:         ":8008",
		ReadTimeout:  h.readTimeout,
		WriteTimeout: h.writeTimeout,
	}
	//h.registrator.RegisterServiceNode(context.Background(), *common.ServiceName, *common.NodeName, *common.IPPort, time.Duration(config.GetConfig().Etcd.TTL)*time.Second)
	//This can actually be made to run in a goroutine
	if err := h.e.StartServer(s); err != nil && err != http.ErrServerClosed {
		h.e.Logger.Error(err)
		h.e.Logger.Info("Shutting down the server")
		os.Exit(1)
	}
}

//etcdctl get --prefix d --user="root" --password="PASSWORD"
