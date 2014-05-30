package rack

import "net/http"

// see wrapper directory to get an idea, how to write a wrapper
// see example directory to know how to use it

type (
	Wrapper interface {
		Wrap(http.Handler) http.Handler
	}

	Racker interface {
		http.Handler
		Wrap(wrapper ...Wrapper)
	}

	rack struct{ http.Handler }
)

// add middleware in the order in which it is
// processed (i.e. reverse of rack.Wrap)
func New(h http.Handler, middleware ...Wrapper) (ø *rack) {
	ø = &rack{h}
	for i := len(middleware) - 1; i >= 0; i-- {
		ø.Handler = middleware[i].Wrap(ø.Handler)
	}
	return
}

// wraps prevoius handlers / wrappers
// the last wrapper is the first in the call
// chain of a request
func (ø *rack) Wrap(wrapper ...Wrapper) {
	for _, wr := range wrapper {
		ø.Handler = wr.Wrap(ø.Handler)
	}
}
