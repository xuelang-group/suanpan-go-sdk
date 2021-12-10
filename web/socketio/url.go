package socketio

import (
	"log"
	"net/url"
)

type URLOptions struct {
	Scheme string
	Path   string
	Query  url.Values
	Logger *log.Logger
}

type URLOption func(o *URLOptions)

func WithPath(p string) URLOption {
	return func(o *URLOptions) {
		o.Path = p
	}
}

func WithScheme(s string) URLOption {
	return func(o *URLOptions) {
		o.Scheme = s
	}
}

func GetURL(host string, opts ...URLOption) *url.URL {
	query := make(url.Values)
	query.Set("EIO", "4")
	query.Set("transport", "websocket")

	options := &URLOptions{
		Scheme: "ws",
		Path:   "/socket.io/",
		Query:  query,
	}

	for _, opt := range opts {
		opt(options)
	}

	u := new(url.URL)
	u.Host = host
	u.Scheme = options.Scheme
	u.Path = options.Path
	u.RawQuery = options.Query.Encode()

	return u
}
