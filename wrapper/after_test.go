package wrapper

import (
	. "launchpad.net/gocheck"
)

type afterSuite struct{}

var _ = Suite(&afterSuite{})

func (s *afterSuite) TestAfter(c *C) {
	r := makeRack(After{webwrite("AFTER")})
	rw, req := newTestRequest("GET", "/")
	r.ServeHTTP(rw, req)
	assertResponse(c, rw, "*AFTER", 200)
}
