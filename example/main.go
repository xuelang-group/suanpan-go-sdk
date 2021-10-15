package main

import (
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/stream"
)

func output(m stream.Message) {
	// m.Data = {
	//     "hello": "world"
	// }
	m.Send(map[string]interface{}{
		"out1": m.Data,
	})
}

func main() {
	msgs := stream.Subscribe()

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			output(msg)
		}
	}()

	<-forever
}