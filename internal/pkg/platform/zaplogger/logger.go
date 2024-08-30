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
	if _, err := os.Stat(fmt.Sprintf("../../logs/%s", *common.ServiceName)); os.IsNotExist(err) {
		os.MkdirAll(fmt.Sprintf("../../logs/%s", *common.ServiceName), 0700)
	}
	logger, err := getFileLogger(fmt.Sprintf("../../logs/%s/app.log", *common.ServiceName))
	if err != nil {
		return nil, err
	}
	return &Logger{
		logger: logger,
	}, nil
}

func (l *Logger) NewInfo(logString string) {
	l.logger.Sugar().Info(logString)
	l.logger.Sync()
}

func (l *Logger) NewError(logString string) {
	l.logger.Sugar().Error(logString)
	l.logger.Sync()
}
