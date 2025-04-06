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
	if _, err := os.Stat("./logs"); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("file exists, error is: %s", err)
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

// Inserts a new info log with new line added to the end into the log file, using zap's sugared logger.
func (l *Logger) NewInfo(logString string) {
	l.logger.Sugar().Info(logString + "\n")
	l.logger.Sync()
}

// Inserts a new error log with new line added to the end into the log file, using zap's sugared logger.
func (l *Logger) NewError(logString string) {
	l.logger.Sugar().Error(logString + "\n")
	l.logger.Sync()
}
