package suanpan_log

import "testing"

func TestLogKitLogger(t *testing.T) {
	logkitLogger := GetLogkitLogger()
	logkitLogger.Info("hi")
}
