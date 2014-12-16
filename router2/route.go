package router2

import (
	"bytes"
	"fmt"
	"gopkg.in/metakeule/meta.v5"
	"gopkg.in/metakeule/rack.v5"
	"gopkg.in/metakeule/rack.v5/helper"
	"net/http"
	"reflect"
	// "net/url"
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
	//Router  *MountedRouter
}

type mountedRoute struct {
	Path   string
	Router *MountedRouter
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

/*
	TODO:

	- make MountedRouter a Racker (fullfill the rack.Racker interface)
	- return routes and allow creation of urls based on routes
*/

type MountedRouter struct {
	*PathNode
	path     string
	parent   *MountedRouter
	routes   map[*route]mountedRoute
	NotFound http.Handler
	wrapper  []rack.Wrapper
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

type Vars struct {
	http.ResponseWriter
	v map[string]string
}

func (v *Vars) Get(key string) string {
	return v.v[key]
}

func (v *Vars) SetStruct(ptrToStruct interface{}, key string) {
	fn := func(f reflect.StructField, fv reflect.Value) {
		tag := f.Tag.Get(key)
		varkey := f.Name
		if tag == "-" {
			return
		}
		if tag != "" {
			varkey = tag
		}
		vv, has := v.v[varkey]
		if has {
			//			meta.Convert(i, t)
			meta.Struct.Set(ptrToStruct, f.Name, vv)
		}
	}
	meta.Struct.EachRaw(ptrToStruct, fn)
}

func (v *Vars) Has(key string) bool {
	_, has := v.v[key]
	return has
}

func (ø *MountedRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	leaf, wc := ø.PathNode.Match(r.URL.Path)
	if leaf != nil {
		h = leaf.Handler(r.Method)
	}

	if h == nil {
		fake := helper.NewFake(w)
		ø.Serve404(fake, r)
		//w.WriteHeader(404)
		// fake.HasChanged()

		if fake.WHeader != 0 {
			w.WriteHeader(fake.WHeader)
		} else {
			w.WriteHeader(404)
		}
		w.Write(fake.Buffer.Bytes())
		return
	}

	//&Vars{w,wc}
	h.ServeHTTP(&Vars{w, wc}, r)
}

func (ø *MountedRouter) Path() string {
	if ø.parent == nil {
		return ø.path
	}
	return path.Join(ø.parent.Path(), ø.path)
}

func (ø *MountedRouter) MustURLMap(rt *route, params map[string]string) string {
	u, err := ø.URLMap(rt, params)
	if err != nil {
		panic(err.Error())
	}
	return u
}

// params are key/values
func (ø *MountedRouter) URLMap(rt *route, params map[string]string) (string, error) {
	mounted, found := ø.routes[rt]
	if !found {
		return "", fmt.Errorf("can't find mounted route for route %#v", rt)
	}

	segments := splitPath(mounted.Path)

	for i := range segments {
		wc, wcName := isWildcard(segments[i])
		if wc {
			repl, ok := params[wcName]
			if !ok {
				return "", fmt.Errorf("missing parameter for %s", wcName)
			}
			segments[i] = repl
		}
	}

	return "/" + strings.Join(segments, "/"), nil
}

var strTy = reflect.TypeOf("")

func (ø *MountedRouter) URLStruct(rt *route, paramStruct interface{}, tagKey string) (string, error) {
	val := reflect.ValueOf(paramStruct)
	if val.Kind() != reflect.Struct {
		panic(fmt.Errorf("%T is not a struct", paramStruct))
	}

	params := map[string]string{}

	fn := func(field reflect.StructField, val reflect.Value, tagVal string) {
		params[tagVal] = val.Convert(strTy).String()
	}
	meta.Struct.EachTag(paramStruct, tagKey, fn)

	return ø.URLMap(rt, params)
}

func (ø *MountedRouter) MustURLStruct(rt *route, paramStruct interface{}, tagKey string) string {
	u, err := ø.URLStruct(rt, paramStruct, tagKey)
	if err != nil {
		panic(err.Error())
	}
	return u
}

// params are key/value pairs
func (ø *MountedRouter) URL(rt *route, params ...string) (string, error) {
	if len(params)%2 != 0 {
		panic("number of params must be even (pairs of key, value)")
	}
	vars := map[string]string{}
	for i := 0; i < len(params)/2; i += 2 {
		vars[params[i]] = params[i+1]
	}

	return ø.URLMap(rt, vars)
}

func (ø *MountedRouter) MustURL(rt *route, params ...string) string {
	u, err := ø.URL(rt, params...)
	if err != nil {
		panic(err.Error())
	}
	return u
}

func (ø *MountedRouter) registerRouteX(p string, rt *route, origRt *route) error {
	if ø.parent == nil {
		pp := path.Join(ø.Path(), p)
		for v, h := range rt.handler {
			if _, isSub := h.(*Router); !isSub {
				r := rack.New(h)
				r.Wrap(ø.wrapper...)
				err := ø.PathNode.addX(pp, v, r)
				if err != nil {
					return fmt.Errorf("can't register route for path %s, verb %s: %s\n", pp, v, err.Error())
				}
			}
		}

		//fmt.Printf("registering route %p\n", origRt)
		/*
			old, has := ø.routes[origRt]
			if has {
				return fmt.Errorf("route %#v already registered at path %#v", origRt, old.Path)
			}
		*/

		mrt := mountedRoute{}
		mrt.Path = pp
		mrt.Router = ø
		ø.routes[origRt] = mrt
		return nil
	}
	pp := path.Join(ø.path, p)
	newrt := NewRoute()
	//newrt.Router = ø

	for v, h := range rt.handler {
		if _, isSub := h.(*Router); !isSub {
			r := rack.New(h)
			r.Wrap(ø.wrapper...)

			err := newrt.AddHandlerX(r, v)
			if err != nil {
				return fmt.Errorf("can't register route for path %s, verb %s: %s\n", pp, v, err.Error())
			}

			//rt.handler[v] = r
		}
	}

	return ø.parent.registerRouteX(pp, newrt, origRt)
	//return ø.parent.registerRouteX(pp, rt)
}

type Router struct {
	wrapper   []rack.Wrapper
	NotFound  http.Handler
	routes    map[string]*route
	subrouter map[string]*Router
}

func New(wrapper ...rack.Wrapper) (ø *Router) {
	ø = &Router{
		wrapper:   wrapper,
		routes:    map[string]*route{},
		subrouter: map[string]*Router{},
	}
	return
}

func newMountedRouter() *MountedRouter {
	return &MountedRouter{
		routes: map[*route]mountedRoute{},
	}
}

func (ø *Router) mountX(mr *MountedRouter, parent *MountedRouter, path string) error {
	mr.path = path
	mr.NotFound = ø.NotFound
	mr.wrapper = ø.wrapper
	mr.parent = parent
	if parent == nil {
		mr.PathNode = newPathNode()
	}
	for p, rt := range ø.routes {
		//fmt.Printf("registering route for path: %s\n", p)
		err := mr.registerRouteX(p, rt, rt)
		if err != nil {
			return err
		}
		for _, h := range rt.handler {
			sub, isSub := h.(*Router)
			if isSub {
				chmr := newMountedRouter()
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
	mr := newMountedRouter()
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

func (ø *Router) Handle(path string, v verb, handler http.Handler) (*route, error) {
	rt, exists := ø.routes[path]
	if exists && rt.Handler(v.String()) != nil {
		panic(fmt.Sprintf("handler for %s (%s) already exists", path, v))
	}
	if !exists {
		rt = NewRoute()
	}
	err := rt.AddHandlerX(handler, v)
	ø.routes[path] = rt
	return rt, err
}

func (r *Router) MustHandle(path string, v verb, handler http.Handler) *route {
	rt, err := r.Handle(path, v, handler)
	if err != nil {
		panic(err.Error())
	}
	return rt
}

func (r *Router) GET(path string, handler http.Handler) *route {
	return r.MustHandle(path, GET, handler)
}
func (r *Router) POST(path string, handler http.Handler) *route {
	return r.MustHandle(path, POST, handler)
}
func (r *Router) PUT(path string, handler http.Handler) *route {
	return r.MustHandle(path, PUT, handler)
}
func (r *Router) DELETE(path string, handler http.Handler) *route {
	return r.MustHandle(path, DELETE, handler)
}
func (r *Router) PATCH(path string, handler http.Handler) *route {
	return r.MustHandle(path, PATCH, handler)
}
func (r *Router) OPTIONS(path string, handler http.Handler) *route {
	return r.MustHandle(path, OPTIONS, handler)
}

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
