package engineio

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/xuelang-group/suanpan-go-sdk/web/socketio/engineio/frame"
	"github.com/xuelang-group/suanpan-go-sdk/web/socketio/engineio/packet"
	"github.com/xuelang-group/suanpan-go-sdk/web/socketio/engineio/session"
	"github.com/xuelang-group/suanpan-go-sdk/web/socketio/engineio/transport"
)

// Opener is client connection which need receive open message first.
type Opener interface {
	Open() (transport.ConnParameters, error)
}

type client struct {
	conn      transport.Conn
	params    transport.ConnParameters
	transport string
	context   interface{}
	close     chan struct{}
	closeOnce sync.Once
}

func (c *client) SetContext(v interface{}) {
	c.context = v
}

func (c *client) Context() interface{} {
	return c.context
}

func (c *client) ID() string {
	return c.params.SID
}

func (c *client) Transport() string {
	return c.transport
}

func (c *client) Close() error {
	c.closeOnce.Do(func() {
		close(c.close)
	})
	return c.conn.Close()
}

func (c *client) NextReader() (session.FrameType, io.ReadCloser, error) {
	for {
		ft, pt, r, err := c.conn.NextReader()
		if err != nil {
			return 0, nil, err
		}

		switch pt {
		case packet.PING:
			logrus.Info("client receive ping....")
			if err = r.Close(); err != nil {
				log.Panic("close reader:", err)
			}

			if err = c.conn.SetReadDeadline(time.Now().Add(c.params.PingInterval + c.params.PingTimeout)); err != nil {
				return 0, nil, err
			}

			logrus.Info("client response pong....")
			w, err := c.conn.NextWriter(frame.String, packet.PONG)
			if err != nil {
				log.Panic(err)
			}
			if err = w.Close(); err != nil {
				log.Panic("close writer:", err)
			}

			if err = c.conn.SetWriteDeadline(time.Now().Add(c.params.PingInterval + c.params.PingTimeout)); err != nil {
				logrus.Error("set writer deadline:", err)
			}

		case packet.CLOSE:
			if err = c.Close(); err != nil {
				log.Panic("close client with packet close:", err)
			}
			return 0, nil, io.EOF

		case packet.MESSAGE:
			return session.FrameType(ft), r, nil
		}
	}
}

func (c *client) NextWriter(typ session.FrameType) (io.WriteCloser, error) {
	return c.conn.NextWriter(frame.Type(typ), packet.MESSAGE)
}

func (c *client) URL() url.URL {
	return c.conn.URL()
}

func (c *client) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *client) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *client) RemoteHeader() http.Header {
	return c.conn.RemoteHeader()
}

func (c *client) serve() {
}
