package logkit

import (
	"net/url"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/util"
	"github.com/xuelang-group/suanpan-go-sdk/web"
	"github.com/xuelang-group/suanpan-go-sdk/web/socketio"
)

var (
	sio *socketio.Conn
	mu sync.Mutex
)

func GetSio() *socketio.Conn {
	mu.Lock()
	defer mu.Unlock()

	if !sio.IsConnected() {
		s, err := buildSio()
		if err == nil {
			sio = s
		}
	}

	return sio
}

func buildSio() (*socketio.Conn, error) {
	e := config.GetEnv()
	u, err := url.Parse(e.SpLogkitUri)
	if err != nil {
		glog.Errorf("parse url error: %w", err)
		return nil, err
	}
	schemeOpt := socketio.WithScheme("ws")
	if u.Scheme == "https" {
		schemeOpt = socketio.WithScheme("wss")
	}
	pathOpt := socketio.WithPath(e.SpLogkitPath)
	socketio.GetURL(u.Host, schemeOpt, pathOpt)

	headerOpt := socketio.WithHeader(web.GetHeaders())
	namespaceOPt := socketio.WithNamespace(e.SpLogkitNamespace)

	u = socketio.GetURL(u.Host, schemeOpt, pathOpt)

	return socketio.New(u.String(), headerOpt, namespaceOPt)
}

func EmitEventLog(title string, level LogLevel)  {
	sio.NextReader()
	//client := web.GetSocketioClient()
	//client.Emit(config.GetEnv().SpLogkitEventsAppend, buildEventLog(title, level))
}

func buildEventLog(title string, level LogLevel) EventLog {
	return EventLog{
		Title: title,
		Level: level.String(),
		Time:  util.ISOString(time.Now()),
		Data: Data{
			Node: config.GetEnv().SpNodeId,
		},
	}
}