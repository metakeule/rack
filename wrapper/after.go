package wrapper

import "net/http"

type (
	// does something after the request has been handled
	After struct{ http.Handler }
)

func (ø After) Wrap(in http.Handler) (out http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in.ServeHTTP(w, r)
		ø.Handler.ServeHTTP(w, r)
	})
}
