package main

import (
	"net/http"
	_ "net/http/pprof"
	// "encoding/json"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
)

func handle(r stream.Request) {
	log.Info("receive request data")
	// m := make(map[string]string)
	// _ = json.Unmarshal([]byte(r.Data), &m)
	// m["hello"] = "world"
	// b, _ := json.Marshal(m)
	// r.Data = string(b)

	r.Send(map[string]string{
		"out1": r.Data,
	})
}

func main() {
	reqs := stream.Subscribe()

	forever := make(chan struct{})

	go func() {
		for req := range reqs {
			handle(req)
		}
	}()

	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	<-forever
}