package router2

import (
	"bytes"
	"fmt"
	"github.com/metakeule/rack"
	"net/http"
	"path"
	"strings"
)

type verb int

const (
	POST verb = 1 << iota
	GET
	PUT
	DELETE
	PATCH
	OPTIONS
)

type route struct {
	handler map[verb]http.Handler
}

func NewRoute() *route {
	r := &route{}
	r.handler = map[verb]http.Handler{}
	return r
}

func (r *route) AddHandlerX(handler http.Handler, v verb) error {
	_, has := r.handler[v]
	if has {
		return fmt.Errorf("handler for verb %s already defined", v)
	}
	r.handler[v] = handler
	return nil
}

func (r *route) Inspect(indent int) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s", strings.Repeat("\t", indent))
	for v, _ := range r.handler {
		fmt.Fprintf(&buf, "%s ", v)
	}
	return buf.String()
}

func (r *route) Handler(v string) http.Handler {
	v = strings.TrimSpace(strings.ToUpper(v))
	ver, ok := verbsFromStrings[v]
	if !ok {
		return nil
	}
	h, exists := r.handler[ver]
	if !exists {
		for vb, h := range r.handler {
			if vb&ver != 0 {
				return h
			}
		}
		return nil
	}
	return h
}

func (v verb) String() string {
	s, exists := verbsToStrings[v]
	if exists {
		return s
	}

	var buf bytes.Buffer
	for vb, ss := range verbsToStrings {
		if vb&v != 0 {
			fmt.Fprintf(&buf, "%s ", ss)
		}
	}
	return buf.String()
}

var verbsToStrings = map[verb]string{
	GET:     "GET",
	POST:    "POST",
	PUT:     "PUT",
	DELETE:  "DELETE",
	PATCH:   "PATCH",
	OPTIONS: "OPTIONS",
}

var verbsFromStrings = map[string]verb{
	"GET":     GET,
	"POST":    POST,
	"PUT":     PUT,
	"DELETE":  DELETE,
	"PATCH":   PATCH,
	"OPTIONS": OPTIONS,
}

type MountedRouter struct {
	*PathNode
	path        string
	parent      *MountedRouter
	NotFound    http.Handler
	middlewares []rack.Wrapper
}

func (ø *MountedRouter) Serve404(w http.ResponseWriter, r *http.Request) {
	if ø.NotFound != nil {
		ø.NotFound.ServeHTTP(w, r)
		return
	}
	if ø.parent != nil {
		ø.parent.Serve404(w, r)
	}
}

func (ø *MountedRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	leaf, _ := ø.PathNode.Match(r.URL.Path)
	if leaf != nil {
		h = leaf.Handler(r.Method)
	}

	if h == nil {
		w.WriteHeader(404)
		ø.Serve404(w, r)
		return
	}
	h.ServeHTTP(w, r)
}

func (ø *MountedRouter) Path() string {
	if ø.parent == nil {
		return ø.path
	}
	return path.Join(ø.parent.Path(), ø.path)
}

func (ø *MountedRouter) registerRouteX(p string, rt *route) error {
	if ø.parent == nil {
		pp := path.Join(ø.Path(), p)
		for v, h := range rt.handler {
			if _, isSub := h.(*Router); !isSub {
				r := rack.New(h)
				r.Wrap(ø.middlewares...)
				err := ø.PathNode.addX(pp, v, r)
				if err != nil {
					return fmt.Errorf("can't register route for path %s, verb %s: %s\n", pp, v, err.Error())
				}
			}
		}
		return nil
	}
	pp := path.Join(ø.path, p)
	newrt := NewRoute()
	for v, h := range rt.handler {
		if _, isSub := h.(*Router); !isSub {
			r := rack.New(h)
			r.Wrap(ø.middlewares...)
			err := newrt.AddHandlerX(r, v)
			if err != nil {
				return fmt.Errorf("can't register route for path %s, verb %s: %s\n", pp, v, err.Error())
			}
		}
	}
	return ø.parent.registerRouteX(pp, newrt)
}

type Router struct {
	middlewares []rack.Wrapper
	NotFound    http.Handler
	routes      map[string]*route
	subrouter   map[string]*Router
}

func New(middlewares ...rack.Wrapper) (ø *Router) {
	ø = &Router{
		middlewares: middlewares,
		routes:      map[string]*route{},
		subrouter:   map[string]*Router{},
	}
	return
}

func (ø *Router) mountX(mr *MountedRouter, parent *MountedRouter, path string) error {
	mr.path = path
	mr.NotFound = ø.NotFound
	mr.middlewares = ø.middlewares
	mr.parent = parent
	if parent == nil {
		mr.PathNode = newPathNode()
	}
	for p, rt := range ø.routes {
		//fmt.Printf("registering route for path: %s\n", p)
		err := mr.registerRouteX(p, rt)
		if err != nil {
			return err
		}
		for _, h := range rt.handler {
			sub, isSub := h.(*Router)
			if isSub {
				chmr := &MountedRouter{}
				e := sub.mountX(chmr, mr, p)
				if e != nil {
					return e
				}
			}
		}
	}
	return nil
}

// method implemented just to fullfill the http.Handler interface, should not be called
// router must be transformed to a MountedRouter via Mount() method in order to serve http
func (ø *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("method implemented just to fullfill the http.Handler interface, should not be called")
}

func (ø *Router) Mount(path string, m *http.ServeMux) (*MountedRouter, error) {
	// fmt.Println("-----------------------------------")
	mr := &MountedRouter{}
	err := ø.mountX(mr, nil, path)
	if err != nil {
		return nil, err
	}
	m.Handle(mr.Path()+"/", mr)
	return mr, nil
}

func (r *Router) MustMount(path string, m *http.ServeMux) *MountedRouter {
	mr, err := r.Mount(path, m)
	if err != nil {
		panic(err.Error())
	}
	return mr
}

func (ø *Router) Handle(path string, v verb, handler http.Handler) error {
	rt, exists := ø.routes[path]
	if exists && rt.Handler(v.String()) != nil {
		panic(fmt.Sprintf("handler for %s (%s) already exists", path, v))
	}
	if !exists {
		rt = NewRoute()
	}
	err := rt.AddHandlerX(handler, v)
	ø.routes[path] = rt
	return err
}

func (r *Router) MustHandle(path string, v verb, handler http.Handler) {
	err := r.Handle(path, v, handler)
	if err != nil {
		panic(err.Error())
	}
}

func (r *Router) GET(path string, handler http.Handler)     { r.MustHandle(path, GET, handler) }
func (r *Router) POST(path string, handler http.Handler)    { r.MustHandle(path, POST, handler) }
func (r *Router) PUT(path string, handler http.Handler)     { r.MustHandle(path, PUT, handler) }
func (r *Router) DELETE(path string, handler http.Handler)  { r.MustHandle(path, DELETE, handler) }
func (r *Router) PATCH(path string, handler http.Handler)   { r.MustHandle(path, PATCH, handler) }
func (r *Router) OPTIONS(path string, handler http.Handler) { r.MustHandle(path, OPTIONS, handler) }

func Mount(path string, r *Router) (*MountedRouter, error) {
	return r.Mount(path, http.DefaultServeMux)
}

func MustMount(path string, r *Router) *MountedRouter {
	mr, err := Mount(path, r)
	if err != nil {
		panic(err.Error())
	}
	return mr
}
