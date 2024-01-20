package logrus

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

const defaultSkip = 2

var defaultLogger *logrusLogger

func init() {
	defaultLogger = NewLogrusLogger()
}

func GetLogger() *logrusLogger {
	return defaultLogger
}

func NewLogrusLogger() *logrusLogger {
	l := new(logrusLogger)
	l.logger = logrus.New()
	l.logger.SetFormatter(
		&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	absPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	l.currentPath = absPath

	l.logger.SetOutput(newOutput())
	l.logger.SetLevel(getLogLevel())
	l.logger.AddHook(newLogrusHook())

	return l
}

type logrusLogger struct {
	logger      *logrus.Logger
	currentPath string
}

func (l *logrusLogger) newEntry() *logrus.Entry {
	_, file, line, ok := runtime.Caller(defaultSkip)
	if ok {
		return l.logger.WithFields(logrus.Fields{
			"location": fmt.Sprintf("%s:%d", path.Base(file), line),
		})
	}

	return l.logger.WithFields(logrus.Fields{})
}

// Debugf implements Logger.
func (l *logrusLogger) Debugf(format string, v ...interface{}) {
	l.newEntry().Debugf(format, v...)
}

// Errorf implements Logger.
func (l *logrusLogger) Errorf(format string, v ...interface{}) {
	l.newEntry().Errorf(format, v...)
}

// Fatalf implements Logger.
func (l *logrusLogger) Fatalf(format string, v ...interface{}) {
	l.newEntry().Fatalf(format, v...)
}

// Infof implements Logger.
func (l *logrusLogger) Infof(format string, v ...interface{}) {
	l.newEntry().Infof(format, v...)
}

// Tracef implements Logger.
func (l *logrusLogger) Tracef(format string, v ...interface{}) {
	l.newEntry().Tracef(format, v...)
}

// Warnf implements Logger.
func (l *logrusLogger) Warnf(format string, v ...interface{}) {
	l.newEntry().Warnf(format, v...)
}

// Debug implements Logger.
func (l *logrusLogger) Debug(v ...interface{}) {
	l.newEntry().Debug(v...)
}

// Error implements Logger.
func (l *logrusLogger) Error(v ...interface{}) {
	l.newEntry().Error(v...)
}

// Fatal implements Logger.
func (l *logrusLogger) Fatal(v ...interface{}) {
	l.newEntry().Fatal(v...)
}

// Info implements Logger.
func (l *logrusLogger) Info(v ...interface{}) {
	l.newEntry().Info(v...)
}

// Trace implements Logger.
func (l *logrusLogger) Trace(v ...interface{}) {
	l.newEntry().Trace(v...)
}

// Warn implements Logger.
func (l *logrusLogger) Warn(v ...interface{}) {
	l.newEntry().Warn(v...)
}

// CtxDebug implements Logger.
func (l *logrusLogger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	l.newEntry().WithContext(ctx).Debugf(format, v...)
}

// CtxError implements Logger.
func (l *logrusLogger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	l.newEntry().WithContext(ctx).Errorf(format, v...)
}

// CtxFatal implements Logger.
func (l *logrusLogger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	l.newEntry().WithContext(ctx).Fatalf(format, v...)
}

// CtxInfo implements Logger.
func (l *logrusLogger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	l.newEntry().WithContext(ctx).Infof(format, v...)
}

// CtxTrace implements Logger.
func (l *logrusLogger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	l.newEntry().WithContext(ctx).Tracef(format, v...)
}

// CtxWarn implements Logger.
func (l *logrusLogger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	l.newEntry().WithContext(ctx).Warnf(format, v...)
}
