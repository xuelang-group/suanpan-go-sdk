package parameter

import "github.com/xuelang-group/suanpan-go-sdk/config"

func Get(param string) string {
	return config.GetArgs()[param]
}