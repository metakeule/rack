package wrapper

import (
	"fmt"
	// "github.com/metakeule/rack"
	// "github.com/metakeule/rack/wrapper"
	"github.com/metakeule/rack"
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
	request, _ := http.NewRequest(method, path, nil)
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

func h(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`*`))
}

func makeRack(wrapper ...rack.Wrapper) rack.Racker {
	r := rack.New(http.HandlerFunc(h))
	r.Wrap(wrapper...)
	return r
}
