package main

import (
	// "encoding/json"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/stream"
)

func handle(r stream.Request) {
	// m := make(map[string]string)
	// _ = json.Unmarshal([]byte(r.Data.(string)), &m)
	// m["hello"] = "world"
	// r.Data, _ = json.Marshal(m)
	r.Send(map[string]interface{}{
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

	<-forever
}