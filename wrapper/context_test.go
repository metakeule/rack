package wrapper

import (
	"github.com/metakeule/rack"
	. "launchpad.net/gocheck"
	"net/http"
)

type contextSuite struct{}

var _ = Suite(&contextSuite{})

type ctx struct {
	path string
	http.ResponseWriter
}

func (c *ctx) Prepare(w http.ResponseWriter, r *http.Request) {
	c.path = r.URL.Path
}

func (c *ctx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("~" + c.path + "~"))
}

func check(w http.ResponseWriter, r *http.Request) {
	c := w.(*ctx)
	w.Write([]byte("#" + c.path + "#"))
}

func (s *contextSuite) TestContextHandlerMethod(c *C) {
	r := rack.New(HandlerMethod((*ctx).ServeHTTP))
	r.Wrap(After{http.HandlerFunc(check)})
	r.Wrap(Before{HandlerMethod((*ctx).Prepare)})
	r.Wrap(Context(&ctx{}))

	rw, req := newTestRequest("GET", "/path")
	r.ServeHTTP(rw, req)
	assertResponse(c, rw, "~/path~#/path#", 200)
}
