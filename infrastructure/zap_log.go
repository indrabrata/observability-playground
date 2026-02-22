package infrastructure

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/indrabrata/observability-playground/constant"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}

func NewZapLog(ctx context.Context) {
	env := os.Getenv("ENVIRONMENT")
	var zConfig zap.Config
	switch env {
	case "PRODUCTION":
		zConfig = zap.NewProductionConfig()
	case "DEVELOPMENT":
		zConfig = zap.NewDevelopmentConfig()
	default:
		zConfig = zap.NewDevelopmentConfig()
	}

	switch os.Getenv("LOG_LEVEL") {
	case "INFO":
		zConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "DEBUG":
		zConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "ERROR":
		zConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "WARN":
		zConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	default:
		zConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	fileName := fmt.Sprintf("./logs/%s.log", time.Now().Format("02-01-2006"))
	ll := lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    1024, //MB
		MaxBackups: 30,
		MaxAge:     90, //days
		Compress:   true,
	}

	zap.RegisterSink("lumberjack", func(*url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: &ll,
		}, nil
	})

	zConfig.Encoding = "json"
	zConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zConfig.OutputPaths = []string{"stdout", fmt.Sprintf("lumberjack:%s", fileName)}

	z, err := zConfig.Build()
	if err != nil {
		panic(err)
	}

	zap.New(otelzap.NewCore(constant.APP_PACKAGE))
	zap.ReplaceGlobals(z)
}
