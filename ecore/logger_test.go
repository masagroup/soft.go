package ecore

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newTestLogger(name string, t require.TestingT) *zap.Logger {
	var logger *zap.Logger
	testLoggerEnv := os.Getenv("TestLogger")
	isLogger, _ := strconv.ParseBool(testLoggerEnv)
	if isLogger {
		f, _ := os.Create(fmt.Sprintf("%s.log", strings.ReplaceAll(name, "/", "_")))
		pe := zap.NewDevelopmentEncoderConfig()
		pe.EncodeTime = zapcore.ISO8601TimeEncoder
		fileEncoder := zapcore.NewJSONEncoder(pe)
		consoleEncoder := zapcore.NewConsoleEncoder(pe)
		level := zap.DebugLevel
		core := zapcore.NewTee(
			zapcore.NewCore(fileEncoder, zapcore.AddSync(f), level),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
		)
		logger = zap.New(core)
	} else {
		logger = zap.NewNop()
		require.NotNil(t, logger)
	}
	return logger
}

type NoWriter struct {
}

func (nw *NoWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (nw *NoWriter) Sync() error {
	return nil
}

func newDebugLogger() *zap.Logger {
	pe := zap.NewDevelopmentEncoderConfig()
	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(pe)
	level := zap.DebugLevel
	core := zapcore.NewTee(zapcore.NewCore(fileEncoder, &NoWriter{}, level))
	return zap.New(core)
}
