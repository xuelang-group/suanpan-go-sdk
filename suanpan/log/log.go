package log

import "github.com/xuelang-group/suanpan-go-sdk/logkit"

func Trace(title string) {
	log(title, logkit.TRACE)
}

func Debug(title string) {
	log(title, logkit.DEBUG)
}

func Info(title string) {
	log(title, logkit.INFO)
}

func Warn(title string) {
	log(title, logkit.WARN)
}

func Error(title string) {
	log(title, logkit.ERROR)
}

func log(title string, level logkit.LogLevel) {
	logkit.EmitEventLog(title, level)
}
