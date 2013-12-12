package wrapper

import (
	"github.com/metakeule/rack"
	. "launchpad.net/gocheck"
	"net/http"
)

type deferSuite struct{}

var _ = Suite(&deferSuite{})

func anyway(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`anyway`))
}

func (s *deferSuite) TestDefer(c *C) {
	r := rack.New(http.HandlerFunc(panicker))
	r.Wrap(Defer{http.HandlerFunc(anyway)})
	rw, req := newTestRequest("GET", "/")
	defer func() { recover() }()
	r.ServeHTTP(rw, req)
	assertResponse(c, rw, "anyway", 200)
}
