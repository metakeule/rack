package wrapper

import (
	. "launchpad.net/gocheck"
)

type aroundSuite struct{}

var _ = Suite(&aroundSuite{})

func (s *aroundSuite) TestAround(c *C) {
	r := makeRack(Around{webwrite("BEFORE"), webwrite("AFTER")})
	rw, req := newTestRequest("GET", "/")
	r.ServeHTTP(rw, req)
	assertResponse(c, rw, "BEFORE*AFTER", 200)
}
