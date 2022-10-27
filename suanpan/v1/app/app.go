package app

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
)

func Run(f func(r stream.Request)) {
	reqs := stream.Subscribe()

	go func() {
		for req := range reqs {
			go f(req)
		}
	}()

	done := make(chan struct{}, 1)

	path := "/internal/trap"
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"success": "true"}`))
			done <- struct{}{}
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"success": "false", "msg": "invalid request"}`))
		}
	})

	go http.ListenAndServe(":" + config.GetEnv().SpTermPort, nil)

	<- done
	log.Info("Exsited")
}