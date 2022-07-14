package app

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
)

func Run(f func(r stream.Request)) {
	reqs := stream.Subscribe()

	forever := make(chan struct{})

	go func() {
		for req := range reqs {
			f(req)
		}
	}()

	go http.ListenAndServe(":6060", nil)

	<-forever
}