package router2

import (
	"fmt"
	// "github.com/metakeule/rack"
	// "github.com/metakeule/rack/wrapper"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type webwrite string

func (ww webwrite) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, string(ww))
}

//
// gocheck: hook into "go test"
//
func Test(t *testing.T) { TestingT(t) }

// Make a testing request
func newTestRequest(method, path string) (*httptest.ResponseRecorder, *http.Request) {
	request, err := http.NewRequest(method, path, nil)
	if err != nil {
		fmt.Printf("could not make request %s (%s): %s\n", path, method, err.Error())
	}
	recorder := httptest.NewRecorder()

	return recorder, request
}

func assertResponse(c *C, rr *httptest.ResponseRecorder, body string, code int) {
	c.Assert(strings.TrimSpace(string(rr.Body.Bytes())), Equals, body)
	c.Assert(rr.Code, Equals, code)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`not found`))
}

func mount(r *Router, mountpoint string) *MountedRouter {
	return r.MustMount(mountpoint, http.NewServeMux())
}
