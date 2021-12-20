package log

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/logkit"
)

var logLevel logkit.LogLevel

func init() {
	env := config.GetEnv()
	switch env.SpLogkitLogsLevel {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
		setLogLevel(logkit.TRACE)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		setLogLevel(logkit.DEBUG)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
		setLogLevel(logkit.INFO)
	case "warning", "warn":
		logrus.SetLevel(logrus.WarnLevel)
		setLogLevel(logkit.WARN)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
		setLogLevel(logkit.ERROR)
	default:
		logrus.Errorf("unknown SpLogkitLogsLevel: %s", env.SpLogkitLogsLevel)
		setLogLevel(logkit.INFO)
	}

	b, err := strconv.ParseBool(env.SpDebug)
	if err != nil {
		logrus.Error(b)
		b = false
	}
	if b {
		logrus.SetLevel(logrus.DebugLevel)
		setLogLevel(logkit.DEBUG)
	}
}

func Trace(title string) {
	log(title, logkit.TRACE)
	logrus.Trace(title)
}

func Tracef(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args)
	Trace(s)
}

func Debug(title string) {
	log(title, logkit.DEBUG)
	logrus.Debug(title)
}

func Debugf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args)
	Debug(s)
}

func Info(title string) {
	log(title, logkit.INFO)
	logrus.Info(title)
}

func Infof(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args)
	logrus.Info(s)
}

func Warn(title string) {
	log(title, logkit.WARN)
	logrus.Warn(title)
}

func Warnf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args)
	Warn(s)
}

func Error(title string) {
	log(title, logkit.ERROR)
	logrus.Error(title)
}

func Errorf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args)
	Error(s)
}

func setLogLevel(level logkit.LogLevel)  {
	logLevel = level
}

func log(title string, level logkit.LogLevel) {
	if level >= logLevel {
		logkit.EmitEventLog(title, level)
	}
}
