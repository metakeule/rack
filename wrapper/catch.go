package wrapper

import (
	"github.com/metakeule/rack/helper"
	"net/http"
)

type (
	// catches an error while executing the request
	Catch func(p interface{}, w http.ResponseWriter, r *http.Request)
)

func (ø Catch) Wrap(in http.Handler) (out http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fake := helper.NewFake(w)
		defer func() {
			if p := recover(); p != nil {
				// delete the headers that have been set
				headers := w.Header()
				for key, _ := range headers {
					headers.Del(key)
				}
				// the func does the rest
				ø(p, w, r)
			} else {
				if fake.WHeader != 0 {
					w.WriteHeader(fake.WHeader)
				}
				w.Write(fake.Buffer.Bytes())
			}
		}()
		in.ServeHTTP(fake, r)
	})
}
