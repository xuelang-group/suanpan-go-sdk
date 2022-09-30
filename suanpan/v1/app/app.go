package app

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
)

func Run(f func(r stream.Request)) {
	reqs := stream.Subscribe()

	// :6060/debug/pprof
	go http.ListenAndServe(":6060", nil)

	for req := range reqs {
		go f(req)
	}
}