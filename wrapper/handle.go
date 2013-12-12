package wrapper

import (
	"net/http"
)

type (
	// casts the Responsewriter to http.Handler in order to write to itself
	Handle struct{}
)

func (ø Handle) Wrap(in http.Handler) (out http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.(http.Handler).ServeHTTP(w, r)
	})
}
