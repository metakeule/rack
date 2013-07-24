package wrapper

import (
	"net/http"
)

type (
	// does something before and after the request is handled further
	Around struct{ Before, After http.Handler }
)

func (ø Around) Wrap(in http.Handler) (out http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ø.Before.ServeHTTP(w, r)
		in.ServeHTTP(w, r)
		ø.After.ServeHTTP(w, r)
	})
}
