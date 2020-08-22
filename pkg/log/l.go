package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var atom zap.AtomicLevel

func InitLogger() func() {
	config := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewJSONEncoder(config)
	atom = zap.NewAtomicLevel()
	logger := zap.New(zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), atom))
	return zap.ReplaceGlobals(logger)
}

func SetLevel(lvl string) {
	atom.SetLevel(parseLevel(lvl))
}

func parseLevel(lvl string) zapcore.Level {
	switch lvl {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.ErrorLevel
	}
}
