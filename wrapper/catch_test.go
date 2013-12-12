package wrapper

import (
	// "fmt"
	"github.com/metakeule/rack"
	. "launchpad.net/gocheck"
	"net/http"
)

type catchSuite struct{}

var _ = Suite(&catchSuite{})

func catcher(p interface{}, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(p.(string)))
}

func panicker(w http.ResponseWriter, r *http.Request) {
	panic("don't panic")
}

func (s *catchSuite) TestDefer(c *C) {
	r := rack.New(http.HandlerFunc(panicker))
	r.Wrap(Catch(catcher))
	rw, req := newTestRequest("GET", "/")
	r.ServeHTTP(rw, req)
	assertResponse(c, rw, "don't panic", 200)
}
