package wrapper

import (
	"gopkg.in/metakeule/rack.v5/helper"
	"net/http"
)

// first will try all given handler until
// the first one returns something
type first []http.Handler

func (f first) Wrap(in http.Handler) (out http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fake := helper.NewFake(w)
		for _, h := range f {
			h.ServeHTTP(fake, r)
			if fake.HasChanged() {
				if fake.WHeader != 0 {
					w.WriteHeader(fake.WHeader)
				}
				w.Write(fake.Buffer.Bytes())
				return
			}
		}
		in.ServeHTTP(w, r)
	})
}

func First(handler ...http.Handler) first {
	return first(handler)
}
