package logkit

import (
	"net/url"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/util"
	"github.com/xuelang-group/suanpan-go-sdk/web"
	"github.com/xuelang-group/suanpan-go-sdk/web/socketio"
)

func getSio() (*socketio.Conn, error) {
	e := config.GetEnv()
	u, err := url.Parse(e.SpLogkitUri)
	if err != nil {
		glog.Errorf("Parse url error: %w", err)
		return nil, err
	}
	schemeOpt := socketio.WithScheme("ws")
	if u.Scheme == "https" {
		schemeOpt = socketio.WithScheme("wss")
	}
	path := e.SpLogkitPath
	if !strings.HasSuffix(path, `/`) {
		path = path + `/`
	}
	pathOpt := socketio.WithPath(path)
	socketio.GetURL(u.Host, schemeOpt, pathOpt)

	headerOpt := socketio.WithHeader(web.GetHeaders())
	namespaceOpt := socketio.WithNamespace(e.SpLogkitNamespace)

	u = socketio.GetURL(u.Host, schemeOpt, pathOpt)

	return socketio.New(u.String(), headerOpt, namespaceOpt)
}

func EmitEventLog(title string, level LogLevel) {
	sio, err := getSio()
	if err != nil {
		glog.Errorf("Get sio error: %w", err)
	}
	sio.Emit(buildEventLog(title, level))
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