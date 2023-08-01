package socketio

import (
	"fmt"

	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xuelang-group/suanpan-go-sdk/web/socketio/engineio"
	"github.com/xuelang-group/suanpan-go-sdk/web/socketio/engineio/transport"
	"github.com/xuelang-group/suanpan-go-sdk/web/socketio/engineio/transport/websocket"
	"github.com/xuelang-group/suanpan-go-sdk/web/socketio/parser"
)

type ClientOptions struct {
	Namespace      string
	Header         http.Header
	Transports     []transport.Transport
	Reconnect      bool
	EventBufferMax int
}

func (c *ClientOptions) getTransport() []transport.Transport {
	if c != nil && len(c.Transports) != 0 {
		return c.Transports
	}
	return []transport.Transport{
		// polling.Default,
		websocket.Default,
	}
}

func (c *ClientOptions) getNamespace() string {
	if c.Namespace == "" {
		return aliasRootNamespace
	}
	return c.Namespace
}

func (c *ClientOptions) getEventBufferMax() int {
	if c.EventBufferMax > 0 {
		return c.EventBufferMax
	}

	return 1000
}

type ClientConn struct {
	// 连接信息
	urlStr    string
	namespace string
	opt       *ClientOptions

	engineio.Conn

	reconnect   bool // 是否断线重连
	isConnected bool // 连接状态

	id uint64

	encoder *parser.Encoder
	decoder *parser.Decoder

	writeChanMax int // 消息缓存区大小
	writeChan    chan parser.Payload
	errorChan    chan error
	quitChan     chan struct{}

	closeOnce sync.Once
	ack       sync.Map
}

func NewClientConn(urlStr string, opt *ClientOptions) (*ClientConn, error) {
	conn := &ClientConn{
		urlStr:    urlStr,
		namespace: opt.getNamespace(),
		opt:       opt,
		reconnect: opt.Reconnect,

		writeChanMax: opt.getEventBufferMax(),
		writeChan:    make(chan parser.Payload, opt.getEventBufferMax()),
	}

	if err := conn.Reconnect(); err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *ClientConn) Close() error {
	var err error

	c.closeOnce.Do(func() {
		err = c.Conn.Close()
		close(c.quitChan)
		c.isConnected = false
		if c.reconnect {
			go func() {
				ticker := time.NewTicker(5 * time.Second)
				for range ticker.C {
					if err := c.Reconnect(); err != nil {
						continue
					}
					return
				}
			}()
		}
	})

	return err
}

func (c *ClientConn) IsConnected() bool {
	return c.isConnected
}

func (c *ClientConn) Reconnect() error {
	if c.isConnected {
		return nil
	}

	Dialer := &engineio.Dialer{
		Transports: c.opt.getTransport(),
	}
	engineConn, err := Dialer.Dial(c.urlStr, c.opt.Header)
	if err != nil {
		return err
	}

	c.Conn = engineConn
	c.namespace = c.opt.getNamespace()
	c.encoder = parser.NewEncoder(engineConn)
	c.decoder = parser.NewDecoder(engineConn)
	c.errorChan = make(chan error)
	c.quitChan = make(chan struct{})

	// At the beginning of a Socket.IO session, the client MUST send a CONNECT packet
	err = c.connectNamespace(c.namespace)
	if err != nil {
		logrus.Errorf("socketio client connect namespace: %v, err: %v", c.namespace, err)
		return err
	}

	c.isConnected = true

	go c.serveWrite()
	go c.serveRead()

	return nil
}

// CLIENT                                                      SERVER

// │  ───────────────────────────────────────────────────────►  │
// │             { type: CONNECT, namespace: "/" }              │
// │  ◄───────────────────────────────────────────────────────  │
// │   { type: CONNECT, namespace: "/", data: { sid: "..." } }  │
func (c *ClientConn) connectNamespace(nsp string) error {
	header := parser.Header{
		Type:      parser.Connect,
		Namespace: nsp,
	}

	if err := c.encoder.Encode(header); err != nil {
		return err
	}

	var event string
	if err := c.decoder.DecodeHeader(&header, &event); err != nil {
		// c.onError(rootNamespace, err)
		return err
	}

	if header.Type != parser.Connect {
		return fmt.Errorf("client conn received packet: %v, but except: %v", header.Type, parser.Connect)
	}

	if err := c.decoder.Close(); err != nil {
		return err
	}

	return nil
}

// low level
func (c *ClientConn) Write(header parser.Header, args ...reflect.Value) {
	data := make([]interface{}, len(args))

	for i := range data {
		data[i] = args[i].Interface()
	}

	pkg := parser.Payload{
		Header: header,
		Data:   data,
	}

	if len(c.writeChan) >= cap(c.writeChan) {
		<-c.writeChan
	}

	c.writeChan <- pkg
}

// up level
func (c *ClientConn) Emit(eventName string, v ...interface{}) {
	header := parser.Header{
		Type: parser.Event,
	}

	if c.namespace != aliasRootNamespace {
		header.Namespace = c.namespace
	}

	if l := len(v); l > 0 {
		last := v[l-1]
		lastV := reflect.TypeOf(last)

		if lastV.Kind() == reflect.Func {
			f := newAckFunc(last)

			header.ID = c.nextID()
			header.NeedAck = true

			c.ack.Store(header.ID, f)
			v = v[:l-1]
		}
	}

	args := make([]reflect.Value, len(v)+1)
	args[0] = reflect.ValueOf(eventName)

	for i := 1; i < len(args); i++ {
		args[i] = reflect.ValueOf(v[i-1])
	}

	c.Write(header, args...)
}

// must run in one sinle goroutine, it would be exited when conn close
func (c *ClientConn) serveWrite() {
	defer func() {
		logrus.Warn("client conn serveWrite goroutine exit.")
		if err := c.Close(); err != nil {
			logrus.Error("close connect:", err)
		}
	}()

	for {
		select {
		case <-c.quitChan:
			return
		case pkg := <-c.writeChan:
			if err := c.encoder.Encode(pkg.Header, pkg.Data); err != nil {
				logrus.Error("client conn ecode error:", err)
				// c.onError(pkg.Header.Namespace, err)
				return
			}
		}
	}
}

// TODO: only working ping/pong on underlying conn, socket event read not support now.
// must run in one sinle goroutine, it would be exited when conn close
func (c *ClientConn) serveRead() {
	defer func() {
		logrus.Warn("client conn serveRead goroutine exit.")
		if err := c.Close(); err != nil {
			logrus.Error("close connect:", err)
		}
	}()

	var event string

	for {
		var header parser.Header

		if err := c.decoder.DecodeHeader(&header, &event); err != nil {
			// c.onError(rootNamespace, err)
			return
		}

		if err := c.decoder.Close(); err != nil {
			return
		}

		if header.Namespace == aliasRootNamespace {
			header.Namespace = rootNamespace
		}

		// logrus.Info(header)

		var err error
		switch header.Type {
		case parser.Ack, parser.Connect, parser.Disconnect:
			logrus.Infof("header type: %v, coming...", header.Type)
		// 	handler, ok := readHandlerMapping[header.Type]
		// 	if !ok {
		// 		return
		// 	}

		// 	err = handler(c, header)
		case parser.Event:
			logrus.Infof("event: %v, coming...", event)
			// err = eventPacketHandler(c, event, header)
		}

		if err != nil {
			logrus.Error("serve read:", err)
			return
		}
	}
}

func (c *ClientConn) nextID() uint64 {
	c.id++

	return c.id
}
