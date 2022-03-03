package main

import (
	// "encoding/json"
	"net/http"
	_ "net/http/pprof"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
)

func handle(r stream.Request) {
	log.Info("receive request data")
	outputData := r.InputData(1)
	// m := make(map[string]string)
	// _ = json.Unmarshal([]byte(outputData), &m)
	// m["hello"] = "world"
	// b, _ := json.Marshal(m)
	// outputData = string(b)

	r.Send(map[string]string{
		"out1": outputData,
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