package zaplogger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getFileWriteSyncer(filePath string) (zapcore.WriteSyncer, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return zapcore.AddSync(file), nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getFileLogger(filePath string) (*zap.Logger, error) {
	writeSyncer, err := getFileWriteSyncer(filePath)
	if err != nil {
		return nil, err
	}
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	return zap.New(core, zap.AddCaller()), nil
}
