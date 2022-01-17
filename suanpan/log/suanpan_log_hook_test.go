package suanpan_log

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSuanpanLog(t *testing.T) {
	log := logrus.New()
	testArr := make([]string, 0)
	hook, err := NewSuanpanLogWithFunc("logkit", logrus.InfoLevel, false, func(entry *logrus.Entry, hook *SuanpanLogHook) error {
		testArr = append(testArr, entry.Message)
		return nil
	})

	assert.NoError(t, err)

	log.Hooks.Add(hook)

	for _, level := range hook.Levels() {
		if len(log.Hooks[level]) != 1 {
			t.Errorf("hook was not added. The length of log.Hooks[%v]: %v", level, len(log.Hooks[level]))
		}
	}

	//error will log
	log.WithField("logkit", true).Error("123")
	assert.Equal(t, testArr[0], "123")

	//not log
	log.WithField("logkit", false).Error("1")
	assert.Equal(t, len(testArr), 1)

	//not log
	log.Error("1")
	assert.Equal(t, len(testArr), 1)

	//not log
	log.Debug("123")
	assert.Equal(t, len(testArr), 1)

	for l := range log.Hooks {
		delete(log.Hooks, l)
	}

	testArr2 := make([]string, 0)
	hook, err = NewSuanpanLogWithFunc("logkit", logrus.InfoLevel, true, func(entry *logrus.Entry, hook *SuanpanLogHook) error {
		testArr2 = append(testArr2, entry.Message)
		return nil
	})
	log.Hooks.Add(hook)
	assert.NoError(t, err)

	//error will log
	log.WithField("logkit", true).Error("123")
	assert.Equal(t, testArr2[0], "123")

	//log
	log.WithField("logkit", false).Error("1")
	assert.Equal(t, len(testArr2), 2)
	assert.Equal(t, testArr2[1], "1")

	//log
	log.Info("17")
	assert.Equal(t, len(testArr2), 3)
	assert.Equal(t, testArr2[2], "17")

	//nog log
	log.Debug("18")
	assert.Equal(t, len(testArr2), 3)
}
