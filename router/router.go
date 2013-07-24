package router

import (
	"github.com/gorilla/mux"
	"github.com/metakeule/rack"
	"net/http"
	"path"
)

type Router struct {
	*mux.Router
	rack        rack.Racker
	path        string
	middlewares []rack.Wrapper
}

func New(path string, middlewares ...rack.Wrapper) (ø *Router) {
	router := mux.NewRouter()
	if path != "" {
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

func (ø *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var m mux.RouteMatch
	if ø.Router.Match(r, &m) {
		ø.rack.ServeHTTP(w, r)
	} else {
		ø.Router.NotFoundHandler.ServeHTTP(w, r)
	}
}
func (ø *Router) Mount(mux *http.ServeMux) { mux.Handle(ø.path+"/", ø) }

func (ø *Router) SubRouter(p string, middlewares ...rack.Wrapper) (rr *Router) {
	middlewares = append(middlewares, ø.middlewares...)
	rr = New(path.Join(ø.path, p), middlewares...)
	rr.Router.NotFoundHandler = ø.Router.NotFoundHandler
	return
}

func Mount(mux *http.ServeMux, r ...*Router) {
	for _, rr := range r {
		rr.Mount(mux)
	}
}
