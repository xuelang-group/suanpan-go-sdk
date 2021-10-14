package config

import (
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/xuelang-group/suanpan-go-sdk/util"
)

var (
	argsMap *map[string]string
	argsOnce sync.Once
)

const ArgPrefix = `--`

func GetArgs() *map[string]string {
	argsOnce.Do(func() {
		argsMap = buildArgs()
	})

	return argsMap
}

func buildArgs() *map[string]string {
	e := GetEnv()
	params, err := util.DecodeBase64(e.SpParam)
	if err != nil {
		glog.Errorf("Decode sp param failed: %v", err)
		return nil
	}

	paramArray := strings.Fields(strings.TrimSpace(params))
	for i := 0; i < len(paramArray); i++ {
		if strings.HasPrefix(paramArray[i], ArgPrefix) {
			if i+1 < len(paramArray) {
				(*argsMap)[paramArray[i]] = paramArray[i+1]
			} else {
				(*argsMap)[paramArray[i]] = ""
			}
		}
	}

	return argsMap
}