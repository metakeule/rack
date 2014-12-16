package wrapper

import (
	"gopkg.in/metakeule/rack.v5/helper"
	"net/http"
)

type (
	// a guard will do something on the responsewriter
	// if further execution is not permitted
	// if it does nothing, the request continues and is handled further down the chain
	Guard struct{ http.Handler }
)

func (ø Guard) Wrap(in http.Handler) (out http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fake := helper.NewFake(w)
		ø.Handler.ServeHTTP(fake, r)
		if fake.HasChanged() {
			if fake.WHeader != 0 {
				w.WriteHeader(fake.WHeader)
			}
			w.Write(fake.Buffer.Bytes())
		} else {
			in.ServeHTTP(w, r)
		}
	})
}

func GuardFunc(fn func(http.ResponseWriter, *http.Request)) Guard {
	return Guard{http.HandlerFunc(fn)}
}
