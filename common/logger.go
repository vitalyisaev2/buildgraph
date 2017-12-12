package common

import "go.uber.org/zap/zapcore"

type Logger interface {
	With(...zapcore.Field) Logger
	Debug(string, ...zapcore.Field)
	Info(string, ...zapcore.Field)
	Warn(string, ...zapcore.Field)
	Error(string, ...zapcore.Field)
	Fatal(string, ...zapcore.Field)
}
