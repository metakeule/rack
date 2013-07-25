package router

import (
	"github.com/gorilla/mux"
	"github.com/metakeule/rack"
	"net/http"
	"path"
)

type Route struct{ *mux.Route }

func (ø *Route) URL(vals ...string) string {
	u, err := ø.Route.URL(vals...)
	if err != nil {
		panic(err.Error())
	}
	return u.String()
}

type Router struct {
	*mux.Router
	rack         rack.Racker
	path         string
	middlewares  []rack.Wrapper
	hasSubroutes bool
}

func New(path string, middlewares ...rack.Wrapper) (ø *Router) {
	router := mux.NewRouter()
	if path != "/" {
		router = router.PathPrefix(path).Subrouter()
	}
	ø = &Router{
		Router:      router,
		path:        path,
		rack:        rack.New(router),
		middlewares: middlewares,
	}
	ø.rack.Wrap(middlewares...)
	return
}

func (ø *Router) Wrap(middlewares ...rack.Wrapper) {
	if ø.hasSubroutes {
		panic("already have subroutes, wrap before subrouting")
	}
	ø.rack.Wrap(middlewares...)
}

func (ø *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var m mux.RouteMatch
	if ø.Router.Match(r, &m) {
		ø.rack.ServeHTTP(w, r)
	} else {
		ø.Router.NotFoundHandler.ServeHTTP(w, r)
	}
}
func (ø *Router) Mount(mux *http.ServeMux) { mux.Handle(ø.path, ø) }

// overwrite the URL method with our own Route struct
func (ø *Router) NewRoute() *Route { return &Route{ø.Router.NewRoute()} }

func (ø *Router) SubRouter(p string, middlewares ...rack.Wrapper) (rr *Router) {
	ø.hasSubroutes = true
	middlewares = append(middlewares, ø.middlewares...)
	p = path.Join(ø.path, p) + "/"
	rr = New(p, middlewares...)
	rr.Router.NotFoundHandler = ø.Router.NotFoundHandler
	return
}

func Mount(mux *http.ServeMux, r ...*Router) {
	for _, rr := range r {
		rr.Mount(mux)
	}
}
