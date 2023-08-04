package logger

import (
	"fmt"
	"github.com/guneyin/bist-tools/pkg/config"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// Interface -.
//type Interface interface {
//	Debug(message interface{}, args ...interface{})
//	Info(message string, args ...interface{})
//	Warn(message string, args ...interface{})
//	Error(message interface{}, args ...interface{})
//	Fatal(message interface{}, args ...interface{})
//}

//var logger *Logger

//type Logger struct {
//	logger *zerolog.Logger
//}

var lg *zerolog.Logger

//var _ Interface = (*Logger)(nil)

func newLogger() *zerolog.Logger {
	var l zerolog.Level

	switch strings.ToLower(config.Cfg.Logger.Level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	skipFrameCount := 3

	zl := zerolog.New(os.Stdout).With().Timestamp().CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).Logger()

	return &zl
}

func Init() {
	lg = newLogger()
}

func Debug(message interface{}, args ...interface{}) {
	msg("debug", message, args...)
}

func Info(message string, args ...interface{}) {
	log(message, args...)
}

func Warn(message string, args ...interface{}) {
	log(message, args...)
}

func Error(message interface{}, args ...interface{}) {
	if lg.GetLevel() == zerolog.DebugLevel {
		Debug(message, args...)
	}

	msg("error", message, args...)
}

func Fatal(message interface{}, args ...interface{}) {
	msg("fatal", message, args...)

	os.Exit(1)
}

func log(message string, args ...interface{}) {
	if len(args) == 0 {
		lg.Info().Msg(message)
	} else {
		lg.Info().Msgf(message, args...)
	}
}

func msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		log(msg.Error(), args...)
	case string:
		log(msg, args...)
	default:
		log(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}
