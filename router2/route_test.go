package router2

import (
	"fmt"
	// "fmt"
	"github.com/metakeule/rack"
	"github.com/metakeule/rack/wrapper"
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

	router := mount(makeRouter(wrapper.Around{webwrite("#"), webwrite("#")}), "/mount")
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

	router := mount(makeRouter(wrapper.Around{webwrite("#"), webwrite("#")}), "/")
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

	inner := makeRouter(wrapper.Around{webwrite("~"), webwrite("~")})
	outer := makeRouter(wrapper.Around{webwrite("#"), webwrite("#")})
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
	inner := New(wrapper.Around{webwrite("~"), webwrite("~")})
	inner.MustHandle("/b", POST, webwrite("B-POST"))

	outer := New(wrapper.Around{webwrite("#"), webwrite("#")})
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
	inner := New(wrapper.Around{webwrite("~"), webwrite("~")})
	inner.MustHandle("/", GET, webwrite("INNER-ROOT"))
	inner.MustHandle("/a", GET|POST, webwrite("A-INNER-GET-POST"))

	outer := New(wrapper.Around{webwrite("#"), webwrite("#")})
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
