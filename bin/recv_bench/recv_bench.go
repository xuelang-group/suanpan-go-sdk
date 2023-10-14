package main

import (
	// "encoding/json"

	"time"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/app"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
)

var recv_pak_cnt uint64 = 0
var startTime = time.Now()
var lastReceivedTime = time.Now()

func handle(r stream.Request) {
	// Check if more than 5 seconds have passed since the last packet
	if time.Since(lastReceivedTime).Seconds() > 5 {
		recv_pak_cnt = 0
		startTime = time.Now()
		log.Infof("---new recv benchmark start, single pak_size:%.3f MB---", float64(len(r.InputData(1)))/(1024*1024))
	}

	lastReceivedTime = time.Now()
	recv_pak_cnt++

	if recv_pak_cnt%1000 == 0 {
		elapsedTime := time.Since(startTime).Seconds()

		// Calculate QPS
		qps := float64(recv_pak_cnt) / elapsedTime

		// Calculate Bandwidth (in MB/s)
		pakSizeMB := float64(len(r.InputData(1))) / (1024 * 1024) // Convert bytes to MB
		bandwidth := (pakSizeMB * float64(recv_pak_cnt)) / elapsedTime

		log.Infof("recv_pak_cnt: %d, QPS: %.2f, Bandwidth: %.2f MB/s",
			recv_pak_cnt, qps, bandwidth)
	}
}

func main() {
	log.Info("start recv bench")
	app.Run(handle)
}
