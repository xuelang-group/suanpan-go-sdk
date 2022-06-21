package desktop

import (
	"time"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/web"
)

func RegisterFreePort(nodePort string) {
	for {
		err := web.RegisterFreePort(nodePort)
		if err != nil {
			log.Warn("retry register free port after 10 seconds")
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}
}

func RegisterPort(nodePort string, port string) {
	for {
		err := web.RegisterPort(nodePort, port)
		if err != nil {
			log.Warn("retry register port after 10 seconds")
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}
}