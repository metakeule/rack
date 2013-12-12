package wrapper

import (
	. "launchpad.net/gocheck"
)

type beforeSuite struct{}

var _ = Suite(&beforeSuite{})

func (s *beforeSuite) TestBefore(c *C) {
	r := makeRack(Before{webwrite("BEFORE")})
	rw, req := newTestRequest("GET", "/")
	r.ServeHTTP(rw, req)
	assertResponse(c, rw, "BEFORE*", 200)
}
