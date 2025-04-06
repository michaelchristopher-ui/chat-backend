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

// NewServer initializes a server instance. Server must be started with StartServer.
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

// GetEcho is a getter for the echo instance
func (h server) GetEcho() *echo.Echo {
	return h.e
}

// StartServer starts the server
func (h server) StartServer() {
	s := &http.Server{
		Addr:         ":8008",
		ReadTimeout:  h.readTimeout,
		WriteTimeout: h.writeTimeout,
	}

	//TODO: To prepare for a more sophisticated service discovery solution, preferably register here with a function with a signature like below.
	// h.registrator.RegisterServiceNode(context.Background(), *common.ServiceName, *common.NodeName, *common.IPPort, time.Duration(config.GetConfig().Etcd.TTL)*time.Second)
	// The function will Record the service, node names and IP for connection purposes, and periodically ping the service discovery server to indicate that the server is still online.

	if err := h.e.StartServer(s); err != nil && err != http.ErrServerClosed {
		h.e.Logger.Error(err)
		h.e.Logger.Info("Shutting down the server")
		os.Exit(1)
	}
}
