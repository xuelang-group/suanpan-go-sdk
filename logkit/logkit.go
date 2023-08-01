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

var (
	sioClientConn  *socketio.ClientConn
	onceWithoutErr util.OnceWithoutErr
)

func initSioClientConn() (*socketio.ClientConn, error) {
	e := config.GetEnv()
	if e.SpLogkitUri == "" {
		return nil, errors.New("SpLogkitUri is empty")
	}
	u, err := url.Parse(e.SpLogkitUri)
	if err != nil {
		logrus.Errorf("Parse url error: %v", err)
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
	u = socketio.GetURL(u.Host, schemeOpt, pathOpt)

	conn, err := socketio.NewClientConn(u.String(), &socketio.ClientOptions{
		Namespace:      e.SpLogkitNamespace,
		Header:         web.GetHeaders(), // not working now
		Reconnect:      true,
		EventBufferMax: 1000,
	})

	return conn, err
}

func getSioClientConn() (*socketio.ClientConn, error) {
	if sioClientConn != nil {
		return sioClientConn, nil
	}

	var err error
	onceWithoutErr.Do(func() error {
		sioClientConn, err = initSioClientConn()
		return err
	})

	if sioClientConn == nil {
		return nil, err
	}

	return sioClientConn, nil
}

func EmitEventLog(title string, level LogLevel) {
	sio, err := getSioClientConn()
	if err != nil {
		logrus.Errorf("Get sio client conn error: %v", err)
		return
	}

	e := buildEvent(title, level)
	sio.Emit(e.Name, e.AppID, e.Log)
}

func buildEvent(title string, level LogLevel) Event {
	return Event{
		Name:  config.GetEnv().SpLogkitEventsAppend,
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
