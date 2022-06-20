package desktop

import (
	"time"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/web"
)

func RegisterPort(nodePort string) {
	for {
		err := web.RegisterPort(nodePort)
		if err != nil {
			log.Warn("retry register port after 10 seconds")
			time.Sleep(10 * time.Second)
		}
	}
}
