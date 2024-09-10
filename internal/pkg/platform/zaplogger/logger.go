package zaplogger

import (
	"fmt"
	"os"
	"websocket_client/internal/common"
	"websocket_client/internal/pkg/core/adapter/loggeradapter"

	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

func NewLogger() (loggeradapter.Adapter, error) {
	if _, err := os.Stat(fmt.Sprintf("./logs")); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("File exists, error is: %s", err)
		}
		err = os.MkdirAll("./logs/", 0777)
		if err != nil {
			return nil, err
		}
	}
	logger, err := getFileLogger("./logs/app.log")
	if err != nil {
		return nil, err
	}
	ret := &Logger{
		logger: logger,
	}
	ret.NewInfo(fmt.Sprintf("logger set up for service %s", *common.ServiceName))
	return ret, nil
}

func (l *Logger) NewInfo(logString string) {
	l.logger.Sugar().Info(logString)
	l.logger.Sync()
}

func (l *Logger) NewError(logString string) {
	l.logger.Sugar().Error(logString)
	l.logger.Sync()
}
