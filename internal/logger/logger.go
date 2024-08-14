package logger

import "go.uber.org/zap"

type Logger interface {
	Info(msg string, err error, opts ...zap.Option)
}

type ZapLogger struct {
	log *zap.Logger
}

func (z *ZapLogger) Info(msg string, err error, fields ...zap.Field) {
	panic("TODO: impl me")
}
