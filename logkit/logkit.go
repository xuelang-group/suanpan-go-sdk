package logkit

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/util"
	"github.com/xuelang-group/suanpan-go-sdk/web"
	"github.com/xuelang-group/suanpan-go-sdk/web/socketio"
)

func getSio() (*socketio.Conn, error) {
	e := config.GetEnv()
	if e.SpLogkitUri == "" {
		return nil, errors.New("SpLogkitUri is empty")
	}
	u, err := url.Parse(e.SpLogkitUri)
	if err != nil {
		logrus.Errorf("Parse url error: %w", err)
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
		logrus.Errorf("Get sio error: %w", err)
		return
	}
	defer sio.Close()

	e := buildEvent(title, level)
	sio.Emit(e.Name, e.AppID, e.Log)
}

func buildEvent(title string, level LogLevel) Event {
	return Event{
		Name: config.GetEnv().SpLogkitEventsAppend,
		AppID: config.GetEnv().SpAppId,
		Log: EventLog{
			Title: title,
			Level: level.String(),
			Time:  util.ISOString(time.Now()),
			Data: Data{
				Node: config.GetEnv().SpNodeId,
			},
		},
	}
}