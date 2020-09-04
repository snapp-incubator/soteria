package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var atom zap.AtomicLevel

// InitLogger will replace global zap.L() with our own preferred configs
func InitLogger() func() {
	config := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewJSONEncoder(config)
	atom = zap.NewAtomicLevel()
	logger := zap.New(zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), atom))
	return zap.ReplaceGlobals(logger)
}

// SetLevel will change zap's log level
func SetLevel(lvl string) {
	atom.SetLevel(parseLevel(lvl))
}

// parseLevel will convert a string based log level to zapcore.Level
func parseLevel(lvl string) zapcore.Level {
	lvl = strings.ToLower(lvl)
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
