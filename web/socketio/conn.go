package socketio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xuelang-group/suanpan-go-sdk/web/socketio/engineio"
	"github.com/xuelang-group/suanpan-go-sdk/web/socketio/ioutil"
)

type Conn struct {
	ws  *websocket.Conn
	wch chan io.Reader

	pingInterval time.Duration
	pingTimeout  time.Duration

	namespace string
}

type ConnOptions struct {
	QueueSize uint
	Header    http.Header
	Dialer    *websocket.Dialer
	Namespace string
}

type ConnOption func(o *ConnOptions)

func WithHeader(h http.Header) ConnOption {
	return func(o *ConnOptions) {
		o.Header = h
	}
}

func WithNamespace(n string) ConnOption {
	return func(o *ConnOptions) {
		o.Namespace = n
	}
}

func New(urlStr string, opts ...ConnOption) (*Conn, error) {
	options := &ConnOptions{
		QueueSize: 100,
		Header:    nil,
		Namespace: "/",
		Dialer: &websocket.Dialer{
			Proxy: http.ProxyFromEnvironment,
		},
	}
	for _, opt := range opts {
		opt(options)
	}

	ws, _, err := options.Dialer.Dial(urlStr, options.Header)
	if err != nil {
		return nil, err
	}

	wch := make(chan io.Reader, options.QueueSize)
	go func() {
		for r := range wch {
			wc, err := ws.NextWriter(websocket.TextMessage)
			if err != nil {
				continue
			}
			if _, err := io.Copy(wc, r); err != nil {
				continue
			}
			wc.Close()
		}
	}()

	mt, r, err := ws.NextReader()
	if err != nil {
		return nil, err
	}
	if mt != websocket.TextMessage {
		return nil, fmt.Errorf("currently supports only text message: %v", mt)
	}

	session, err := engineio.ReadHandshake(r)
	c := &Conn{
		ws:           ws,
		wch:          wch,
		pingInterval: time.Duration(session.PingInterval) * time.Millisecond,
		pingTimeout:  time.Duration(session.PingTimeout) * time.Millisecond,
		namespace:    options.Namespace,
	}

	err = c.ConnectNamespace(options.Namespace)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Conn) ConnectNamespace(namespace string) error {
	p := &Packet{
		Type:      CONNECT,
		Namespace: namespace,
		ID:        -1,
	}
	w := c.NewWriter(engineio.MessagePrefix())
	if err := NewEncoder(w).Encode(p); err != nil {
		return err
	}
	w.Flush()
	return nil
}

func (c *Conn) NextReader() (io.Reader, error) {
	mt, r, err := c.ws.NextReader()
	if err != nil {
		return nil, err
	}
	if mt != websocket.TextMessage {
		return nil, fmt.Errorf("currently supports only text message: %v", mt)
	}
	return r, nil
}

func (c *Conn) NewWriter(p []byte) *ioutil.Writer {
	w := &ioutil.Writer{
		Prefix: p,
		Ch:     c.wch,
		Buf:    &bytes.Buffer{},
	}

	if p != nil {
		w.Write(p)
	}

	return w
}

func (c *Conn) Emit(args ...interface{}) error {
	p := &Packet{
		Type:      EVENT,
		Namespace: c.namespace,
		ID:        -1,
	}
	w := c.NewWriter(engineio.MessagePrefix())
	if err := NewEncoder(w).Encode(p); err != nil {
		return fmt.Errorf("encode header: %w", err)
	}

	b, err := json.Marshal([]interface{}{args})
	if err != nil {
		return err
	}
	w.Write(b)
	w.Flush()
	_, err = c.NextReader()
	return err
}
