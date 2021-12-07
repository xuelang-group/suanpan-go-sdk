package ioutil

import (
	"bytes"
	"io"
)

type Writer struct {
	Prefix []byte
	Ch chan io.Reader
	Buf *bytes.Buffer
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.Buf.Write(p)
}

func (w *Writer) Flush() {
	w.Ch <-w.Buf
}