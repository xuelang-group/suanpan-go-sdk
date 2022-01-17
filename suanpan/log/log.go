package suanpan_log

import (
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/logkit"
)

func init() {
	env := config.GetEnv()
	level := logrus.WarnLevel
	switch env.SpLogkitLogsLevel {
	case "trace":
		level = logrus.TraceLevel
	case "debug":
		level = logrus.DebugLevel
	case "info":
		level = logrus.InfoLevel
	case "warning", "warn":
		level = logrus.WarnLevel
	case "error":
		level = logrus.ErrorLevel
	default:
		logrus.Errorf("unknown SpLogkitLogsLevel: %s", env.SpLogkitLogsLevel)
		level = logrus.InfoLevel
	}

	b, err := strconv.ParseBool(env.SpDebug)
	if err != nil {
		logrus.Error(b)
		b = false
	}
	if b {
		logrus.SetLevel(logrus.DebugLevel)
		level = logrus.DebugLevel
	}

	hook, err := NewSuanpanLogWithFunc("logkit", level, false, func(entry *logrus.Entry, hook *SuanpanLogHook) error {
		level := logkit.INFO
		switch entry.Level {
		case logrus.TraceLevel:
			level = logkit.TRACE
		case logrus.DebugLevel:
			level = logkit.DEBUG
		case logrus.InfoLevel:
			level = logkit.INFO
		case logrus.WarnLevel:
			level = logkit.WARN
		case logrus.ErrorLevel:
			level = logkit.ERROR
		}

		//need async ? go log(msg, level)??
		log(entry.Message, level)
		return nil
	})

	if err != nil {
		logrus.Error(err)
	}
	logrus.AddHook(hook)
}

func log(title string, level logkit.LogLevel) {
	logkit.EmitEventLog(title, level)
}

func GetLogkitLogger() *logrus.Entry {
	return logrus.WithField("logkit", true)
}
