package config

import (
	"encoding/base64"
	"strings"
	"sync"

	"github.com/golang/glog"
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
	params, err := base64.StdEncoding.DecodeString(e.SpParam)
	if err != nil {
		glog.Errorf("Decode sp param failed: %w", err)
		return nil
	}

	paramArray := strings.Fields(strings.TrimSpace(string(params)))
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