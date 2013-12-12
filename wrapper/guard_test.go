package wrapper

import (
	. "launchpad.net/gocheck"
	"net/http"
)

type guardSuite struct{}

var _ = Suite(&guardSuite{})

func (s *guardSuite) TestGuardForbidden(c *C) {
	r := makeRack(Guard{webwrite("forbidden")})
	rw, req := newTestRequest("GET", "/")
	r.ServeHTTP(rw, req)
	assertResponse(c, rw, "forbidden", 200)
}

func doNothing(w http.ResponseWriter, r *http.Request) {}

func (s *guardSuite) TestGuardPassthrough(c *C) {
	r := makeRack(Guard{http.HandlerFunc(doNothing)})
	rw, req := newTestRequest("GET", "/")
	r.ServeHTTP(rw, req)
	assertResponse(c, rw, "*", 200)
}
