package main

import (
	// "encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/stream"
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

	<-forever
}
