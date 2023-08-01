package engineio

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/xuelang-group/suanpan-go-sdk/web/socketio/engineio/packet"
	"github.com/xuelang-group/suanpan-go-sdk/web/socketio/engineio/transport"
)

// Dialer is dialer configure.
type Dialer struct {
	Transports []transport.Transport
}

// Dial returns a connection which dials to url with requestHeader.
func (d *Dialer) Dial(urlStr string, requestHeader http.Header) (Conn, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		logrus.Error("parse url err: ", err)
		return nil, err
	}

	query := u.Query()
	query.Set("EIO", "4")
	u.RawQuery = query.Encode()

	var conn transport.Conn

	for i := len(d.Transports) - 1; i >= 0; i-- {
		if conn != nil {
			if closeErr := conn.Close(); closeErr != nil {
				logrus.Error("close connect:", closeErr)
			}
		}

		t := d.Transports[i]

		conn, err = t.Dial(u, requestHeader)
		if err != nil {
			logrus.Error("transport dial:", err)
			continue
		}

		var params transport.ConnParameters
		if p, ok := conn.(Opener); ok {
			params, err = p.Open()
			if err != nil {
				logrus.Error("open transport connect:", err)
				continue
			}
		} else {
			var pt packet.Type
			var r io.ReadCloser

			_, pt, r, err = conn.NextReader()
			if err != nil {
				continue
			}

			func() {
				defer func() {
					if closeErr := r.Close(); closeErr != nil {
						logrus.Error("close connect reader:", closeErr)
					}
				}()

				if pt != packet.OPEN {
					err = errors.New("invalid open")
					return
				}

				params, err = transport.ReadConnParameters(r)
				if err != nil {
					return
				}
			}()
		}
		if err != nil {
			logrus.Error("transport dialer:", err)
			continue
		}

		// log.Printf("params: %+v", params)
		ret := &client{
			conn:      conn,
			params:    params,
			transport: t.Name(),
			close:     make(chan struct{}),
		}

		// go ret.serve()

		return ret, nil
	}

	return nil, err
}
