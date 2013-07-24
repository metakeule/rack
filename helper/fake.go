package helper

import (
	"bytes"
	"net/http"
)

type (
	fake struct {
		w       http.ResponseWriter
		Buffer  bytes.Buffer
		WHeader int
		changed bool
	}
)

func (ø *fake) HasChanged() bool    { return ø.changed }
func (ø *fake) Header() http.Header { ø.changed = true; return ø.w.Header() }
func (ø *fake) WriteHeader(i int)   { ø.changed = true; ø.WHeader = i }
func (ø *fake) Write(b []byte) (int, error) {
	ø.changed = true
	return ø.Buffer.Write(b)
}

func NewFake(w http.ResponseWriter) (ø *fake) {
	ø = &fake{}
	ø.w = w
	return
}
