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

func New(h http.Handler) (ø *rack) { return &rack{h} }

func (ø *rack) Wrap(wrapper ...Wrapper) {
	for _, wr := range wrapper {
		ø.Handler = wr.Wrap(ø.Handler)
	}
}
