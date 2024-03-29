package worker

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct{}

var _ asynq.Logger = (*Logger)(nil)

func NewLogger() *Logger { return &Logger{} }

func (logger *Logger) Print(logLevel zerolog.Level, args ...interface{}) {
	log.WithLevel(logLevel).Msg(fmt.Sprint(args...))
}

func (logger *Logger) Info(args ...interface{}) {
	logger.Print(zerolog.InfoLevel, args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.Print(zerolog.ErrorLevel, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.Print(zerolog.FatalLevel, args...)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.Print(zerolog.DebugLevel, args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.Print(zerolog.WarnLevel, args...)
}
