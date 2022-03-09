package parameter

import (
	"strconv"

	"github.com/xuelang-group/suanpan-go-sdk/config"
)

const ParamPrefix = `param`

func Get(param string) string {
	return config.GetArgs()[config.ArgNamePrefix + param]
}

func GetParam(i int) string {
	return Get(ParamPrefix + strconv.Itoa(i))
}