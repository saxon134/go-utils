package saLog

import (
	"fmt"
	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"time"
)

type zapLog struct {
	logger *zap.Logger
}

func (m *zapLog) Log(a ...interface{}) {
	m.logger.Info(fmt.Sprint(a))
}

func initZapLog() *zapLog {
	m := new(zapLog)

	fileWriter, err := zaprotatelogs.New(
		path.Join("log", "%Y-%m-%d.log"),
		zaprotatelogs.WithLinkName("log"),
		zaprotatelogs.WithMaxAge(7*24*time.Hour),
		zaprotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return nil
	}

	writer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter))
	config := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	m.logger = zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(config), writer, zapcore.InfoLevel))
	return m
}
