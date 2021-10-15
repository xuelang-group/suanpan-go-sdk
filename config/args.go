package config

import (
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/xuelang-group/suanpan-go-sdk/util"
)

var (
	argsMap = make(map[string]string)
	argsOnce sync.Once
)

const (
	ArgNamePrefix = `--`
	ArgValuePrefix = `'`
)

func GetArgs() map[string]string {
	argsOnce.Do(func() {
		argsMap = buildArgs()
	})

	return argsMap
}

func buildArgs() map[string]string {
	e := GetEnv()
	params, err := util.DecodeBase64(e.SpParam)
	if err != nil {
		glog.Errorf("Decode sp param failed: %v", err)
		return nil
	}

	paramArray := strings.Fields(strings.TrimSpace(params))
	for i := 0; i < len(paramArray); i++ {
		if strings.HasPrefix(paramArray[i], ArgNamePrefix) &&
			i+1 < len(paramArray) &&
			strings.HasPrefix(paramArray[i+1], ArgValuePrefix) {
			if i+1 < len(paramArray) {
				argsMap[paramArray[i]] = strings.Trim(paramArray[i+1], ArgValuePrefix)
				i++
			} else {
				argsMap[paramArray[i]] = ""
			}
		}
	}

	return argsMap
}