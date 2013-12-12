package wrapper

import (
	. "launchpad.net/gocheck"
	"net/http"
)

type firstSuite struct{}

var _ = Suite(&firstSuite{})

func (s *firstSuite) TestFirstA(c *C) {
	r := makeRack(First(webwrite("a"), webwrite("b")))
	rw, req := newTestRequest("GET", "/")
	r.ServeHTTP(rw, req)
	assertResponse(c, rw, "a", 200)
}

func (s *firstSuite) TestFirstB(c *C) {
	r := makeRack(First(http.HandlerFunc(doNothing), webwrite("b")))
	rw, req := newTestRequest("GET", "/")
	r.ServeHTTP(rw, req)
	assertResponse(c, rw, "b", 200)
}

func (s *firstSuite) TestFirstPassthrough(c *C) {
	r := makeRack(First(http.HandlerFunc(doNothing), http.HandlerFunc(doNothing)))
	rw, req := newTestRequest("GET", "/")
	r.ServeHTTP(rw, req)
	assertResponse(c, rw, "*", 200)
}
