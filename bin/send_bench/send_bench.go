package main

import (
	// "encoding/json"

	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
)

func generateStringWithSize(size int) string {
	// Calculate the number of times the character "X" should be repeated
	repeatCount := size / len("X")

	// Create a string with the desired size by repeating the character "X"
	result := strings.Repeat("X", repeatCount)

	return result
}

func main() {
	e := config.GetEnv()
	sp_os := e.SpOs
	//use sp_os to get data_size and cycle_cnt
	sp_os_list := strings.Split(sp_os, "_")
	data_size := 1024 * 1 //100KB
	cycle_cnt := 1000
	send_worker := 2
	if len(sp_os_list) == 3 {
		// Convert the string to an integer
		num, err := strconv.Atoi(sp_os_list[0])
		if err != nil {
			log.Errorf("strconv.Atoi datasize error:%v", err)
		} else {
			data_size = num
		}

		num, err = strconv.Atoi(sp_os_list[1])
		if err != nil {
			log.Errorf("strconv.Atoi cycle_cnt error:%v", err)
		} else {
			cycle_cnt = num
		}

		num, err = strconv.Atoi(sp_os_list[2])
		if err != nil {
			log.Errorf("strconv.Atoi cycle_cnt error:%v", err)
		} else {
			send_worker = num
		}
	}

	// Convert bytes to MB for logging
	log.Infof("start send bench with data size:%.2f MB, cycle cnt:%d, worker:%d", float64(data_size)/(1024*1024), cycle_cnt, send_worker)

	var wg sync.WaitGroup
	total_cnt := 0

	total_cnt_chan := make(chan struct{}, 1000)
	main_donw_ch := make(chan struct{})

	go func() {
		for {
			startTime := time.Now()
			for i := 1; i <= cycle_cnt; i++ {
				<-total_cnt_chan
				total_cnt++

				//print if ok
				if i%1000 == 0 {
					elapsedTime := time.Since(startTime).Seconds()
					bandwidth := float64(i*data_size) / (1024 * 1024 * elapsedTime)
					qps := float64(i) / elapsedTime
					log.Infof("send pak cnt: %d, current send bandwith: %.3f MB/s, current QPS: %.2f pak/s", i, bandwidth, qps)
				}
			}
			finalElapsedTime := time.Since(startTime).Seconds()
			avgBandwidth := float64(total_cnt*data_size) / (1024 * 1024 * finalElapsedTime)
			avgQPS := float64(total_cnt) / finalElapsedTime
			log.Infof("Average Bandwidth: %.3f MB/s, Average QPS: %.2f pak/s", avgBandwidth, avgQPS)

			main_donw_ch <- struct{}{}
		}
	}()

	data := generateStringWithSize(data_size)
	for i := 0; i < send_worker; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()

			//r := stream.Request{ID: "123", Extra: "456"}
			data_map := map[string]string{"out1": data}
			data_map["success"] = "true"
			//data_map["request_id"] = ""
			data_map["node_id"] = config.GetEnv().SpNodeId

			//worker cnt
			worker_cnt := cycle_cnt / send_worker
			logrus.Infof("start worker %d, worker cnt:%d", workerId, worker_cnt)
			for i := 0; i < worker_cnt; i++ {
				//r.SendSuccess(data_map)
				//r.SendOutput(1, data)
				stream.Send(data_map)
				//s := stream.GetStream()
				//s.StreamDirectSend(data_map)
				total_cnt_chan <- struct{}{}
			}
		}(i)
	}
	wg.Wait()
	<-main_donw_ch
	logrus.Info("all jobs complete")
}
