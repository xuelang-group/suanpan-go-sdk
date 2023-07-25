package ioutil

import (
	"bytes"
	"io"
)

type Writer struct {
	Prefix []byte
	Ch     chan io.Reader
	Buf    *bytes.Buffer
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.Buf.Write(p)
}

func (w *Writer) Flush() {
	// Remove one item from chan if full
	if len(w.Ch) == cap(w.Ch) {
		// Channel was full, but might not be by now
		select {
		case <-w.Ch:
		// Discard one item
		default:
			// Maybe it was empty already
		}
	}
	w.Ch <- w.Buf
}
