package wrapper

import "net/http"

type (
	// defer something after the request handling
	Defer struct{ http.Handler }
)

func (ø Defer) Wrap(in http.Handler) (out http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() { ø.ServeHTTP(w, r) }()
		in.ServeHTTP(w, r)
	})
}

func DeferFunc(fn func(http.ResponseWriter, *http.Request)) Defer {
	return Defer{http.HandlerFunc(fn)}
}
