package wrapper

import (
	"github.com/metakeule/rack"
	. "launchpad.net/gocheck"
	"net/http"
)

type handleSuite struct{}

var _ = Suite(&handleSuite{})

type handle struct {
	path string
	http.ResponseWriter
}

func (c *handle) Prepare(w http.ResponseWriter, r *http.Request) {
	c.path = r.URL.Path
}

func (c *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("~" + c.path + "~"))
}

func (s *handleSuite) TestContextHandlerMethod(c *C) {
	r := rack.New(nil)
	r.Wrap(Handle{})
	r.Wrap(Before{HandlerMethod((*handle).Prepare)})
	r.Wrap(Context(handle{}))

	rw, req := newTestRequest("GET", "/path")
	r.ServeHTTP(rw, req)
	assertResponse(c, rw, "~/path~", 200)
}
