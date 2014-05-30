package router2

import (
	"fmt"
	// "fmt"
	"github.com/metakeule/rack"
	wr "github.com/metakeule/rack/wrapper"
	. "launchpad.net/gocheck"
	"net/http"
)

type routetest struct {
	path string
	body string
	code int
}

type routeSuite struct{}

var _ = Suite(&routeSuite{})

func makeRouter(mw ...rack.Wrapper) *Router {
	corpus := []routetest{
		{"/a", "A", 200},
		{"/b", "B", 200},
		{"/x", "", 404},
		{"/a/x", "AX", 200},
		{"/a/b", "AB", 200},
		{"/b/x", "BX", 200},
		{"/:sth/x", "SthX", 200},
	}

	router := New(mw...)
	router.NotFound = http.HandlerFunc(notFound)
	for _, r := range corpus {
		if r.code == 200 {
			router.MustHandle(r.path, GET, webwrite(r.body))
		}
	}
	return router
}

func (s *routeSuite) TestRouting(c *C) {
	corpus := []routetest{
		{"/a", "A", 200},
		{"/b", "B", 200},
		{"/x", "not found", 404},
		{"/a/x", "AX", 200},
		{"/b/x", "BX", 200},
		{"/z/x", "SthX", 200},
		{"/y", "not found", 404},
	}

	router := mount(makeRouter(), "/")

	for _, r := range corpus {
		rw, req := newTestRequest("GET", r.path)
		router.ServeHTTP(rw, req)
		assertResponse(c, rw, r.body, r.code)
	}
}

func (s *routeSuite) TestRoutingMiddlewareMounted(c *C) {
	corpus := []routetest{
		{"/mount/a", "#A#", 200},
		{"/mount/b", "#B#", 200},
		{"/mount/x", "not found", 404},
		{"/mount/a/x", "#AX#", 200},
		{"/mount/b/x", "#BX#", 200},
		{"/mount/z/x", "#SthX#", 200},
		{"/mount/y", "not found", 404},
		{"/a", "not found", 404},
		{"/z/x", "not found", 404},
	}

	router := mount(makeRouter(wr.Around{webwrite("#"), webwrite("#")}), "/mount")
	for _, r := range corpus {
		rw, req := newTestRequest("GET", r.path)
		router.ServeHTTP(rw, req)
		assertResponse(c, rw, r.body, r.code)
	}
}

func (s *routeSuite) TestRoutingMiddleware(c *C) {
	corpus := []routetest{
		{"/a", "#A#", 200},
		{"/b", "#B#", 200},
		{"/x", "not found", 404},
		{"/a/x", "#AX#", 200},
		{"/b/x", "#BX#", 200},
		{"/z/x", "#SthX#", 200},
		{"/y", "not found", 404},
	}

	router := mount(makeRouter(wr.Around{webwrite("#"), webwrite("#")}), "/")
	for _, r := range corpus {
		rw, req := newTestRequest("GET", r.path)
		router.ServeHTTP(rw, req)
		assertResponse(c, rw, r.body, r.code)
	}
}

func (s *routeSuite) TestRoutingMounted(c *C) {
	corpus := []routetest{
		{"/mount/a", "A", 200},
		{"/mount/b", "B", 200},
		{"/mount/x", "not found", 404},
		{"/mount/a/x", "AX", 200},
		{"/mount/b/x", "BX", 200},
		{"/mount/z/x", "SthX", 200},
		{"/mount/y", "not found", 404},
		{"/a", "not found", 404},
		{"/z/x", "not found", 404},
	}

	router := mount(makeRouter(), "/mount")

	for _, r := range corpus {
		rw, req := newTestRequest("GET", r.path)
		router.ServeHTTP(rw, req)
		assertResponse(c, rw, r.body, r.code)
	}
}

func (s *routeSuite) TestRoutingSubroutes(c *C) {
	corpus := []routetest{
		{"/outer/a", "A", 200},
		{"/outer/b", "B", 200},
		{"/outer/x", "not found", 404},
		{"/outer/a/x", "AX", 200},
		{"/outer/b/x", "BX", 200},
		{"/outer/z/x", "SthX", 200},
		{"/outer/y", "not found", 404},
		{"/a", "not found", 404},
		{"/z/x", "not found", 404},

		{"/outer/inner/a", "A", 200},
		{"/outer/inner/b", "B", 200},
		{"/outer/inner/a/x", "AX", 200},
		{"/outer/inner/b/x", "BX", 200},
		{"/outer/inner/z/x", "SthX", 200},
		{"/outer/inner/y", "not found", 404},
		{"/inner/a", "not found", 404},
		{"/inner/z/x", "not found", 404},
	}
	inner := makeRouter()
	outer := makeRouter()
	outer.MustHandle("/inner", GET, inner)

	router := mount(outer, "/outer")
	_ = router
	for _, r := range corpus {
		rw, req := newTestRequest("GET", r.path)
		router.ServeHTTP(rw, req)
		assertResponse(c, rw, r.body, r.code)
	}
}

func (s *routeSuite) TestRoutingMiddlewareSubroutes(c *C) {
	corpus := []routetest{
		{"/outer/a", "#A#", 200},
		{"/outer/b", "#B#", 200},
		{"/outer/x", "not found", 404},
		{"/outer/a/x", "#AX#", 200},
		{"/outer/b/x", "#BX#", 200},
		{"/outer/z/x", "#SthX#", 200},
		{"/outer/y", "not found", 404},
		{"/a", "not found", 404},
		{"/z/x", "not found", 404},

		{"/outer/inner/a", "#~A~#", 200},
		{"/outer/inner/b", "#~B~#", 200},
		{"/outer/inner/a/x", "#~AX~#", 200},
		{"/outer/inner/b/x", "#~BX~#", 200},
		{"/outer/inner/z/x", "#~SthX~#", 200},
		{"/outer/inner/y", "not found", 404},
		{"/inner/a", "not found", 404},
		{"/inner/z/x", "not found", 404},
	}

	inner := makeRouter(wr.Around{webwrite("~"), webwrite("~")})
	outer := makeRouter(wr.Around{webwrite("#"), webwrite("#")})
	outer.MustHandle("/inner", GET, inner)

	router := mount(outer, "/outer")
	//fmt.Println(router.Inspect(0))
	_ = router
	for _, r := range corpus {
		rw, req := newTestRequest("GET", r.path)
		router.ServeHTTP(rw, req)
		assertResponse(c, rw, r.body, r.code)
	}
}

func (s *routeSuite) TestRoutingVerbs(c *C) {
	r := makeRouter()
	r.MustHandle("/a", POST, webwrite("A-POST"))
	router := mount(r, "/")

	rw, req := newTestRequest("GET", "/a")
	router.ServeHTTP(rw, req)
	assertResponse(c, rw, "A", 200)

	rw, req = newTestRequest("POST", "/a")
	router.ServeHTTP(rw, req)
	assertResponse(c, rw, "A-POST", 200)
}

func (s *routeSuite) TestRoutingHandlerAndSubroutes(c *C) {
	inner := New(wr.Around{webwrite("~"), webwrite("~")})
	inner.MustHandle("/b", POST, webwrite("B-POST"))

	outer := New(wr.Around{webwrite("#"), webwrite("#")})
	outer.MustHandle("/a/b", GET, webwrite("B-GET"))
	outer.MustHandle("/a", POST, inner)
	outer.MustHandle("/other", POST, inner)

	//	fmt.Println(outer.Inspect(0))
	router := mount(outer, "/mount")
	// fmt.Println(router.Inspect(0))

	rw, req := newTestRequest("GET", "/mount/a/b")
	router.ServeHTTP(rw, req)
	assertResponse(c, rw, "#B-GET#", 200)

	rw, req = newTestRequest("POST", "/mount/a/b")
	router.ServeHTTP(rw, req)
	assertResponse(c, rw, "#~B-POST~#", 200)

	rw, req = newTestRequest("POST", "/mount/other/b")
	router.ServeHTTP(rw, req)
	assertResponse(c, rw, "#~B-POST~#", 200)
}

func (s *routeSuite) TestRoutingHandlerCombined(c *C) {
	inner := New(wr.Around{webwrite("~"), webwrite("~")})
	inner.MustHandle("/", GET, webwrite("INNER-ROOT"))
	inner.MustHandle("/a", GET|POST, webwrite("A-INNER-GET-POST"))

	outer := New(wr.Around{webwrite("#"), webwrite("#")})
	outer.MustHandle("/a", GET|POST, webwrite("A-OUTER-GET-POST"))

	outer.MustHandle("/inner", GET|POST, inner)

	_ = fmt.Println
	//	fmt.Println(outer.Inspect(0))
	router := mount(outer, "/mount")
	// fmt.Println(router.Inspect(0))
	rw, req := newTestRequest("GET", "/mount/a")
	router.ServeHTTP(rw, req)
	assertResponse(c, rw, "#A-OUTER-GET-POST#", 200)

	rw, req = newTestRequest("POST", "/mount/a")
	router.ServeHTTP(rw, req)
	assertResponse(c, rw, "#A-OUTER-GET-POST#", 200)

	rw, req = newTestRequest("GET", "/mount/inner")
	router.ServeHTTP(rw, req)
	assertResponse(c, rw, "#~INNER-ROOT~#", 200)

	rw, req = newTestRequest("OPTIONS", "/mount/a")
	router.ServeHTTP(rw, req)
	assertResponse(c, rw, "", 404)

	rw, req = newTestRequest("GET", "/mount/inner/a")
	router.ServeHTTP(rw, req)
	assertResponse(c, rw, "#~A-INNER-GET-POST~#", 200)

	rw, req = newTestRequest("POST", "/mount/inner/a")
	router.ServeHTTP(rw, req)
	assertResponse(c, rw, "#~A-INNER-GET-POST~#", 200)

	rw, req = newTestRequest("PATCH", "/mount/inner/a")
	router.ServeHTTP(rw, req)
	assertResponse(c, rw, "", 404)

}

func (s *routeSuite) TestRoutingSubRouteRoot(c *C) {
	admin := New()
	admin.MustHandle("/", GET, webwrite("ADMIN"))
	index := New()
	index.MustHandle("/admin", GET, admin)

	router := mount(index, "/index")

	rw, req := newTestRequest("GET", "/index/admin/")
	router.ServeHTTP(rw, req)
	assertResponse(c, rw, "ADMIN", 200)

	rw, req = newTestRequest("GET", "/index/admin")
	router.ServeHTTP(rw, req)
	assertResponse(c, rw, "ADMIN", 200)
}

type v struct {
	x string
	y string
}

func (vv *v) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	vars := &Vars{}
	wr.UnWrap(rw, &vars)
	vv.x = vars.Get("x")
	vv.y = vars.Get("y")
}

func (s *routeSuite) TestVars(c *C) {
	vv := &v{}
	r := New()
	r.MustHandle("/a/:x/c/:y", GET, vv)
	router := mount(r, "/r")
	rw, req := newTestRequest("GET", "/r/a/b/c/d")
	router.ServeHTTP(rw, req)
	//assertResponse(c, rw, "ADMIN", 200)
	c.Assert(vv.x, Equals, "b")
	c.Assert(vv.y, Equals, "d")
}

type ctx struct {
	App  string `var:"app"`
	path string
	http.ResponseWriter
}

func (c *ctx) SetPath(w http.ResponseWriter, r *http.Request) {
	c.path = r.URL.Path
}

func (c *ctx) SetVars(vars *Vars, w http.ResponseWriter, r *http.Request) {
	vars.SetStruct(c, "var")
}

func (c *ctx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("app: " + c.App + " path: " + c.path))
}

func (s *routeSuite) TestVarsSetStruct(c *C) {
	ct := &ctx{}
	r := New(
		wr.Before{wr.HandlerMethod(ct.SetPath)},
		wr.Before{wr.HandlerMethod(ct.SetVars)},
		wr.Context(ct))
	r.MustHandle("/app/:app/hiho", GET, wr.HandlerMethod(ct.ServeHTTP))

	router := mount(r, "/r")
	rw, req := newTestRequest("GET", "/r/app/X/hiho")
	router.ServeHTTP(rw, req)
	assertResponse(c, rw, "app: X path: /r/app/X/hiho", 200)
}

type uStr1 struct {
	Y string `urltest:"y"`
}

func (s *routeSuite) TestURL(c *C) {
	admin := New()
	route1 := admin.MustHandle("/x", GET, webwrite("ADMIN-X"))
	route2 := admin.MustHandle("/:y/z", GET, webwrite("ADMIN-Z"))
	index1 := New()
	index1.MustHandle("/admin1", GET, admin)
	index2 := New()
	index2.MustHandle("/admin2", GET, admin)

	router1 := mount(index1, "/index1")
	router2 := mount(index2, "/index2")

	url1 := router1.MustURL(route1)
	c.Assert(url1, Equals, "/index1/admin1/x")
	url2 := router2.MustURL(route1)
	c.Assert(url2, Equals, "/index2/admin2/x")

	url3 := router1.MustURL(route2, "y", "p")
	c.Assert(url3, Equals, "/index1/admin1/p/z")

	url4 := router2.MustURL(route2, "y", "p")
	c.Assert(url4, Equals, "/index2/admin2/p/z")

	_, err := router1.URL(route2)
	c.Assert(err, NotNil)

	str1 := uStr1{"q"}
	url5 := router1.MustURLStruct(route2, str1, "urltest")
	c.Assert(url5, Equals, "/index1/admin1/q/z")
}
