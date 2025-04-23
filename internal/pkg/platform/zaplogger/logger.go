package zaplogger

import (
	"fmt"
	"os"
	"runtime"
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

func (l *Logger) getCaller() string {
	/*
		GetCaller is called by Logger functions such as NewInfo, which is then called by the functions that
		require logging. i.e.:
		calling function -> NewXXX -> l.getCaller()

		Hence we would need to ascend two times.
		Should we fail to retrieve the caller, the exact log string would need to be searched.
	*/
	counter, _, _, success := runtime.Caller(2)
	if !success {
		return "unknownCaller"
	}
	return runtime.FuncForPC(counter).Name()
}

// Inserts a new info log with new line added to the end into the log file, using zap's sugared logger.
func (l *Logger) NewInfo(logString string, params ...any) {
	l.logger.Sugar().Infof(logString, l.getCaller(), params)
}

// Inserts a new error log with new line added to the end into the log file, using zap's sugared logger.
func (l *Logger) NewError(logString string, params ...any) {
	l.logger.Sugar().Errorf(logString, l.getCaller(), params)
}
