package logrus

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	logFileName = "./log/%s.log"
	// LogFileMaxSize 每个日志文件最大 MB
	logFileMaxSize = 512
	// LogFileMaxBackups 保留日志文件个数
	logFileMaxBackups = 10
	// LogFileMaxAge 保留日志最大天数
	logFileMaxAge   = 14
	defaultLogLevel = logrus.DebugLevel

	envLogOutputFileName      = "log_output_file_name"
	envLogsSetOutputLocalFile = "log_set_local_file"
	envLogsLevel              = "log_level"
)

func newOutput() io.Writer {
	if os.Getenv(envLogsSetOutputLocalFile) == "true" {
		return io.MultiWriter(
			os.Stdout,
			&lumberjack.Logger{
				Filename:   fmt.Sprintf(logFileName, os.Getenv(envLogOutputFileName)),
				MaxSize:    logFileMaxSize,
				MaxAge:     logFileMaxAge,
				MaxBackups: logFileMaxBackups,
				LocalTime:  true,
				Compress:   false,
			})
	}

	return os.Stdout
}

func getLogLevel() logrus.Level {
	defaultLevel := defaultLogLevel
	switch os.Getenv(envLogsLevel) {
	case "trace":
		defaultLevel = logrus.TraceLevel
	case "debug":
		defaultLevel = logrus.DebugLevel
	case "info":
		defaultLevel = logrus.InfoLevel
	case "warn":
		defaultLevel = logrus.WarnLevel
	case "error":
		defaultLevel = logrus.ErrorLevel
	case "fatal":
		defaultLevel = logrus.FatalLevel
	}

	return defaultLevel
}
