package router2

import (
	"net/http"
	"testing"
)

type routeTest struct {
	route string
	get   string
	vars  map[string]string
}

type handler struct {
	Name string
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func TestRoutes(t *testing.T) {
	table := []routeTest{
		{
			route: "/",
			get:   "/",
			vars:  nil,
		},
		{
			route: "/api/action",
			get:   "/api/action",
			vars:  nil,
		},
		{
			route: "/admin/action",
			get:   "/admin/action",
			vars:  nil,
		},
		{
			route: "/admin/action.json",
			get:   "/admin/action.json",
			vars:  nil,
		},
		{
			route: "/:api/action",
			get:   "/poop/action",
			vars:  map[string]string{"api": "poop"},
		},
		{
			route: "/api/:action",
			get:   "/api/poop",
			vars:  map[string]string{"action": "poop"},
		},
		{
			route: "/:seg1/:seg2/bob",
			get:   "/a/b/bob",
			vars:  map[string]string{"seg1": "a", "seg2": "b"},
		},
		{
			route: "/:seg1/:seg2/ron",
			get:   "/c/d/ron",
			vars:  map[string]string{"seg1": "c", "seg2": "d"},
		},
		{
			route: "/:seg1/:seg2/:seg3",
			get:   "/c/d/wat",
			vars:  map[string]string{"seg1": "c", "seg2": "d", "seg3": "wat"},
		},
		{
			route: "/:seg1/:seg2/ron/apple",
			get:   "/c/d/ron/apple",
			vars:  map[string]string{"seg1": "c", "seg2": "d"},
		},
		{
			route: "/:seg1/:seg2/ron/:apple",
			get:   "/c/d/ron/orange",
			vars:  map[string]string{"seg1": "c", "seg2": "d", "apple": "orange"},
		},
		{
			route: "/site2/:id",
			get:   "/site2/123",
			vars:  map[string]string{"id": "123"},
		},
	}

	n := newPathNode()
	// Create routes
	for _, rt := range table {
		err := n.addX(rt.route, GET, &handler{rt.route})
		if err != nil {
			panic(err.Error())
		}
	}
	for _, rt := range table {
		leaf, wc := n.Match(rt.get)
		if leaf == nil {
			t.Errorf("got no leaf for %#v", rt.route)
		}

		for k, v := range wc {
			exp, has := rt.vars[k]
			if !has {
				t.Errorf("missing key in vars: %s", k)
			}
			if v != exp {
				t.Errorf("incorrect vars value for key %s: %#v, expected: %#v", k, v, exp)
			}
		}
	}
}
