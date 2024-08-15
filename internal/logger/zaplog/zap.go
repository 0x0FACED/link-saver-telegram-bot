package zaplog

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	log *zap.Logger
}

func New() *ZapLogger {
	dirName := "logs"
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		log.Fatalln("cant make dir: ", err)
		return nil
	}

	filename := time.Now().Format("2006-01-02_15-04-05") + ".log"
	filePath := filepath.Join(dirName, filename)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("cant open file: ", err)
		return nil
	}

	config := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	cEnc := zapcore.NewConsoleEncoder(config)
	fEnc := zapcore.NewConsoleEncoder(config)

	core := zapcore.NewTee(
		zapcore.NewCore(cEnc, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		zapcore.NewCore(fEnc, zapcore.AddSync(file), zapcore.InfoLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &ZapLogger{
		log: logger,
	}
}

func (z *ZapLogger) Info(msg string, err error, opts ...zap.Field) {
	panic("not implemented") // TODO: Implement
}

func (z *ZapLogger) Debug(msg string, err error, opts ...zap.Field) {
	panic("not implemented") // TODO: Implement
}

func (z *ZapLogger) Error(msg string, err error, opts ...zap.Field) {
	panic("not implemented") // TODO: Implement
}

func (z *ZapLogger) Fatal(msg string, err error, opts ...zap.Field) {
	panic("not implemented") // TODO: Implement
}
