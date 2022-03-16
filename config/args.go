package config

import (
	"encoding/base64"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
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
		logrus.Errorf("Decode sp param failed: %v", err)
		return nil
	}

	paramArray := strings.Fields(strings.TrimSpace(string(params)))
	addParamArray(paramArray)
	addParamArray(os.Args[1:])

	return argsMap
}

func addParamArray(paramArray []string) {
	for i := 0; i < len(paramArray); i++ {
		if strings.HasPrefix(paramArray[i], ArgNamePrefix) &&
			i+1 < len(paramArray) {
			if i+1 < len(paramArray) {
				argsMap[paramArray[i]] = strings.Trim(paramArray[i+1], ArgValuePrefix)
				i++
			} else {
				argsMap[paramArray[i]] = ""
			}
		}
	}
}