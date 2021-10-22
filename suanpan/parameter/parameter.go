package parameter

import "github.com/xuelang-group/suanpan-go-sdk/config"

func Get() map[string]string {
	return config.GetArgs()
}