package wrapper

import (
	"fmt"
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

type blindctx struct {
	http.ResponseWriter
}

func (c *ctx) Prepare(w http.ResponseWriter, r *http.Request) {
	c.path = r.URL.Path
}

func (c *ctx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("~" + c.path + "~"))
}

func check(w http.ResponseWriter, r *http.Request) {
	_ = rack.New
	c := &ctx{}
	err := UnWrap(w, &c)
	if err != nil {
		panic(err.Error())
	}
	w.Write([]byte("#" + c.path + "#"))
}

func (s *contextSuite) TestContextHandlerMethod(c *C) {
	r := rack.New(HandlerMethod((*ctx).ServeHTTP))
	r.Wrap(After{http.HandlerFunc(check)})
	r.Wrap(Before{HandlerMethod((*ctx).Prepare)})
	r.Wrap(Context(&ctx{}))
	r.Wrap(Context(&blindctx{}))

	rw, req := newTestRequest("GET", "/path")
	r.ServeHTTP(rw, req)
	assertResponse(c, rw, "~/path~#/path#", 200)
}

func (s *contextSuite) TestContextUnwrapIdentical(c *C) {
	c1 := ctx{path: "x"}
	c2 := &ctx{}

	err := UnWrap(&c1, &c2)
	c.Assert(c2.path, Equals, "x")
	if err != nil {
		panic(err.Error())
	}

	/*
		c3 := &ctx{}
		err = UnWrap(&c1, &c3)
		c.Assert(c3.path, Equals, "x")

		if err != nil {
			panic(err.Error())
		}
	*/
}

func (s *contextSuite) TestContextUnwrapNested(c *C) {
	c1 := blindctx{&ctx{path: "x"}}
	c2 := &ctx{}

	err := UnWrap(&c1, &c2)
	if err != nil {
		panic(err.Error())
	}
	c.Assert(c2.path, Equals, "x")
}

func (s *contextSuite) TestContextUnwrapError(c *C) {
	_ = fmt.Println
	c1 := blindctx{}
	c2 := &ctx{}
	err := UnWrap(&c1, &c2)
	// fmt.Println(err.Error())
	c.Assert(err, NotNil)

	rw, _ := newTestRequest("GET", "/path")
	c1 = blindctx{rw}
	err = UnWrap(&c1, &c2)
	// fmt.Println(err.Error())
	c.Assert(err, NotNil)
}
