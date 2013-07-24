package wrapper

import "net/http"

type (
	// does something before the request is handled further
	Before struct{ http.Handler }
)

func (ø Before) Wrap(in http.Handler) (out http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ø.Handler.ServeHTTP(w, r)
		in.ServeHTTP(w, r)
	})
}
