package suanpan_log

import (
	"github.com/sirupsen/logrus"
)

//set handler
type HookFunc func()

type SuanpanLogHook struct {
	HookMsg         string
	IsLogHookAllMsg bool
	levels          []logrus.Level
	fireFunc        fireFunc
}

type fireFunc func(entry *logrus.Entry, hook *SuanpanLogHook) error

//active log hook by default Logkit All msg
//or use logrus field logrus.withfield("loghook", true)
func NewSuanpanLogWithFunc(hookMsg string, level logrus.Level, isLogHookAllMsg bool, fireFunc fireFunc) (*SuanpanLogHook, error) {
	var levels []logrus.Level
	for _, l := range []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	} {
		if l <= level {
			levels = append(levels, l)
		}
	}
	return &SuanpanLogHook{HookMsg: hookMsg, levels: levels, fireFunc: fireFunc, IsLogHookAllMsg: isLogHookAllMsg}, nil
}

func (hook *SuanpanLogHook) Fire(entry *logrus.Entry) error {
	needFire := entry.Data[hook.HookMsg]
	if needFire == true || hook.IsLogHookAllMsg {
		return hook.fireFunc(entry, hook)
	} else {
		//ignore this log msg
		return nil
	}
}

func (hook *SuanpanLogHook) Levels() []logrus.Level {
	return hook.levels
}
